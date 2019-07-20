package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/xuebing1110/notify-inspect/pkg/plugin"
	"github.com/xuebing1110/notify-inspect/pkg/plugin/storage"
	"github.com/xuebing1110/notify-inspect/pkg/schedule"
	"github.com/xuebing1110/notify-inspect/pkg/schedule/cron"
)

func ListPlugins(ctx *gin.Context) {
	SendNormalResponse(ctx, regServer.ListPlugins())
}

func GetPlugin(ctx *gin.Context) {
	pid := ctx.Param("pid")
	pinfo, found := regServer.GetPlugin(pid)
	if !found {
		SendResponse(ctx, http.StatusNotFound, "unknownPlugin", "the plugin was not found")
		return
	}

	SendNormalResponse(ctx, pinfo)
}

func CallPluginSubscribe(ctx *gin.Context) {
	i, found := ctx.Get("subscribe")
	if !found {
		SendResponse(ctx, http.StatusInternalServerError, "ReadSubscribeFromCtxFailed", "read subscribe from context failed")
		return
	}

	subscribe := i.(*plugin.Subscribe)

	p, found := plugin.DefaultRegisterServer.GetPlugin(subscribe.PluginId)
	if !found || p.LostTime > 0 {
		SendResponse(ctx, http.StatusServiceUnavailable, "PluginOffline", fmt.Sprintf("the plugin %s is offline", subscribe.PluginId))
		return
	}

	resp := p.BackendSubscribe(ctx.Request.Context(), subscribe)
	if resp.Code >= 400 {
		ctx.JSON(resp.Code, resp)
		return
	}

	SendNormalResponse(ctx, nil)
}

// user plugin setting
func SavePluginSubscribe(ctx *gin.Context) {

	// subscribe
	subscribe := new(plugin.Subscribe)
	err := ctx.BindJSON(subscribe)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, "ParseSubscribeFailed", err.Error())
		return
	}

	uid := ctx.GetString(CTX_USERID)
	if uid != "" {
		subscribe.UserId = uid
	}

	// save
	err = storage.GlobalStorage.SaveSubscribe(subscribe)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, "SaveSubscribeFailed", err.Error())
		return
	}

	// save to context
	ctx.Set("subscribe", subscribe)
}

func DeletePluginSubscribe(ctx *gin.Context) {
	uid := ctx.GetString(CTX_USERID)
	if uid == "" {
		sendNoUserResponse(ctx)
		return
	}

	pid := ctx.Param("pid")
	if pid == "" {
		SendResponse(ctx, http.StatusInternalServerError, "unknownPluginId", "read pluginid failed")
		return
	}

	err := storage.GlobalStorage.DeleteSubscribe(uid, pid)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, "DeleteSubscribeFailed", err.Error())
		return
	}

	SendNormalResponse(ctx, nil)
}

func GetPluginSubscribe(ctx *gin.Context) {
	uid := ctx.GetString(CTX_USERID)
	if uid == "" {
		sendNoUserResponse(ctx)
		return
	}

	pid := ctx.Param("pid")
	if pid == "" {
		SendResponse(ctx, http.StatusInternalServerError, "unknownPluginId", "read pluginid failed")
		return
	}

	// get configure
	subscribe, err := storage.GlobalStorage.GetSubscribe(uid, pid)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, "GetPluginSubscribeFailed", err.Error())
		return
	}

	// save to context
	ctx.Set("plugin", pid)
	ctx.Set("subscribe", subscribe)
}

func CallPluginSubscribeStatus(ctx *gin.Context) {
	i, found := ctx.Get("subscribe")
	if !found {
		SendResponse(ctx, http.StatusInternalServerError, "ReadSubscribeFromCtxFailed", "read subscribe from context failed")
		return
	}

	subscribe := i.(*plugin.Subscribe)
	pid := ctx.GetString("plugin")

	// get plugin
	p, found := plugin.DefaultRegisterServer.GetPlugin(pid)
	if !found || p.LostTime > 0 {
		subscribe.SetNotAvaliable(fmt.Sprintf("the plugin %s is offline", pid))
		SendNormalResponse(ctx, subscribe)
		// SendResponse(ctx, http.StatusServiceUnavailable, "PluginOffline")
		return
	}

	// get plugin subscribe status
	resp := p.BackendGetSubscribe(ctx.Request.Context(), subscribe)
	if resp.Code < 400 {
		subscribe.SetNotAvaliable(resp.Message)
	}

	SendNormalResponse(ctx, subscribe)
}

// user plugin records
func ListPluginRecords(ctx *gin.Context) {
	uid := ctx.GetString(CTX_USERID)
	if uid == "" {
		sendNoUserResponse(ctx)
		return
	}

	pid := ctx.Param("pid")
	if pid == "" {
		SendResponse(ctx, http.StatusInternalServerError, "unknownPluginId", "read pluginid failed")
		return
	}

	prs, err := storage.GlobalStorage.ListPluginRecords(uid, pid)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, "ListPluginRecordsFailed", err.Error())
		return
	}

	SendNormalResponse(ctx, prs)
}

func GetPluginRecord(ctx *gin.Context) {
	uid := ctx.GetString(CTX_USERID)
	if uid == "" {
		sendNoUserResponse(ctx)
		return
	}

	pid := ctx.Param("pid")
	if pid == "" {
		SendResponse(ctx, http.StatusInternalServerError, "unknownPluginId", "read pluginid failed")
		return
	}

	rid := ctx.Param("rid")
	if rid == "" {
		SendResponse(ctx, http.StatusInternalServerError, "unknownRecordId", "read recordid failed")
		return
	}

	pr, err := storage.GlobalStorage.GetPluginRecord(uid, pid, rid)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, "GetPluginRecordFailed", err.Error())
		return
	}

	SendNormalResponse(ctx, pr)
}

func AddPluginRecord(ctx *gin.Context) {
	uid := ctx.GetString(CTX_USERID)
	if uid == "" {
		sendNoUserResponse(ctx)
		return
	}

	pid := ctx.Param("pid")
	if pid == "" {
		SendResponse(ctx, http.StatusInternalServerError, "unknownPluginId", "read pluginid failed")
		return
	}

	// PluginRecord
	record := new(plugin.PluginRecord)
	err := ctx.BindJSON(record)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, "ParsePluginRecordFailed", err.Error())
		return
	}
	record.UserId = uid
	record.PluginId = pid
	if record.Id == "" {
		record.Id = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	if record.Cron == nil {
		SendResponse(ctx, http.StatusBadRequest, "ParsePluginRecordFailed", "cron is require")
		return
	}

	// put task
	err = schedule.DefaultScheduler.PutTask(record.GetCronTask(), time.Now())
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, "AddSchedulerTaskFailed", err.Error())
		return
	}

	// save record
	err = storage.GlobalStorage.SavePluginRecord(record)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, "SavePluginRecordFailed", err.Error())
		return
	}

	SendNormalResponse(ctx, nil)
}

func ModifyPluginRecord(ctx *gin.Context) {
	uid := ctx.GetString(CTX_USERID)
	if uid == "" {
		sendNoUserResponse(ctx)
		return
	}

	pid := ctx.Param("pid")
	if pid == "" {
		SendResponse(ctx, http.StatusInternalServerError, "unknownPluginId", "read pluginid failed")
		return
	}

	rid := ctx.Param("rid")
	if rid == "" {
		SendResponse(ctx, http.StatusInternalServerError, "unknownRecordId", "read recordid failed")
		return
	}

	// PluginRecord
	pmap := make(map[string]interface{})
	err := ctx.BindJSON(pmap)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, "ParsePluginRecordFailed", err.Error())
		return
	}
	if len(pmap) == 0 {
		SendResponse(ctx, http.StatusBadRequest, "ParsePluginRecordFailed", "body must be not empty")
		return
	}

	// put task
	cron_map, found := pmap["cron"]
	if found {
		cron_setting, err := cron.NewCronTaskSettingFromMap2(cron_map.(map[string]interface{}))
		if err != nil {
			SendResponse(ctx, http.StatusInternalServerError, "ParseCronTaskFailed", err.Error())
			return
		}

		task := &cron.CronTask{plugin.GenerateRecordIdentify(uid, pid, rid), cron_setting}
		err = schedule.DefaultScheduler.PutTask(task, time.Now())
		if err != nil {
			SendResponse(ctx, http.StatusInternalServerError, "AddSchedulerTaskFailed", err.Error())
			return
		}
	}

	// modify record
	found, err = storage.GlobalStorage.ModifyPluginRecord(uid, pid, rid, pmap)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, "ModifyPluginRecord", err.Error())
		return
	}
	if !found {
		SendResponse(ctx, http.StatusNotFound, "PluginRecordNotFound", "can't found the plugin record")
		return
	}

	SendNormalResponse(ctx, nil)
}

func DeletePluginRecord(ctx *gin.Context) {
	uid := ctx.GetString(CTX_USERID)
	if uid == "" {
		return
	}

	pid := ctx.Param("pid")
	if pid == "" {
		SendResponse(ctx, http.StatusInternalServerError, "unknownPluginId", "read pluginid failed")
		return
	}

	rid := ctx.Param("rid")
	if rid == "" {
		SendResponse(ctx, http.StatusInternalServerError, "unknownRecordId", "read recordid failed")
		return
	}

	// remove task
	taskid := plugin.GenerateRecordIdentify(uid, pid, rid)
	err := schedule.DefaultScheduler.RemoveTask(taskid)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, "DelSchedulerTaskFailed", err.Error())
		return
	}

	err = storage.GlobalStorage.DeletePluginRecord(uid, pid, rid)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, "DeletePluginRecord", err.Error())
		return
	}

	SendNormalResponse(ctx, nil)
}

func sendNoUserResponse(ctx *gin.Context) {
	SendResponse(ctx, http.StatusInternalServerError, "unknownUser", "get userid from context failed")
}
