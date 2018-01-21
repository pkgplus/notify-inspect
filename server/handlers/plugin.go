package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kataras/iris/context"

	"github.com/xuebing1110/notify-inspect/pkg/plugin"
	"github.com/xuebing1110/notify-inspect/pkg/plugin/storage"
	"github.com/xuebing1110/notify-inspect/pkg/schedule"
	"github.com/xuebing1110/notify-inspect/pkg/schedule/cron"
)

func ListPlugins(ctx context.Context) {
	SendNormalResponse(ctx, regServer.ListPlugins())
}

func GetPlugin(ctx context.Context) {
	pid := ctx.Params().Get("pid")
	pinfo, found := regServer.GetPlugin(pid)
	if !found {
		SendResponse(ctx, http.StatusNotFound, "unknownPlugin", "the plugin was not found")
		return
	}

	SendNormalResponse(ctx, pinfo)
}

// user plugin setting
func SavePluginSubscribe(ctx context.Context) {
	uid := ctx.Values().GetString(CTX_USERID)
	if uid == "" {
		SendResponse(ctx, http.StatusInternalServerError, "unknownUser", "get userid from context failed")
		return
	}

	// subscribe
	subscribe := new(plugin.Subscribe)
	err := ctx.ReadJSON(subscribe)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, "ParseSubscribeFailed", err.Error())
		return
	}
	subscribe.UserId = uid

	p, found := plugin.DefaultRegisterServer.GetPlugin(subscribe.PluginId)
	if !found {
		SendResponse(ctx, http.StatusServiceUnavailable, "PluginOffline", fmt.Sprintf("the plugin %s is offline", subscribe.PluginId))
		return
	}

	resp := p.BackendSubscribe(subscribe)
	if resp.Code >= 400 {
		ctx.JSON(resp)
		return
	}

	// save
	err = storage.GlobalStorage.SaveSubscribe(subscribe)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, "SaveSubscribeFailed", err.Error())
		return
	}

	SendNormalResponse(ctx, nil)
}

func DeletePluginSubscribe(ctx context.Context) {
	uid := ctx.Values().GetString(CTX_USERID)
	if uid == "" {
		sendNoUserResponse(ctx)
		return
	}

	pid := ctx.Params().Get("pid")
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

func GetPluginSubscribe(ctx context.Context) {
	uid := ctx.Values().GetString(CTX_USERID)
	if uid == "" {
		sendNoUserResponse(ctx)
		return
	}

	pid := ctx.Params().Get("pid")
	if pid == "" {
		SendResponse(ctx, http.StatusInternalServerError, "unknownPluginId", "read pluginid failed")
		return
	}

	ps, err := storage.GlobalStorage.GetSubscribe(uid, pid)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, "GetPluginSubscribeFailed", err.Error())
		return
	}

	SendNormalResponse(ctx, ps)
}

// user plugin records
func ListPluginRecords(ctx context.Context) {
	uid := ctx.Values().GetString(CTX_USERID)
	if uid == "" {
		sendNoUserResponse(ctx)
		return
	}

	pid := ctx.Params().Get("pid")
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

func GetPluginRecord(ctx context.Context) {
	uid := ctx.Values().GetString(CTX_USERID)
	if uid == "" {
		sendNoUserResponse(ctx)
		return
	}

	pid := ctx.Params().Get("pid")
	if pid == "" {
		SendResponse(ctx, http.StatusInternalServerError, "unknownPluginId", "read pluginid failed")
		return
	}

	rid := ctx.Params().Get("rid")
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

func AddPluginRecord(ctx context.Context) {
	uid := ctx.Values().GetString(CTX_USERID)
	if uid == "" {
		sendNoUserResponse(ctx)
		return
	}

	pid := ctx.Params().Get("pid")
	if pid == "" {
		SendResponse(ctx, http.StatusInternalServerError, "unknownPluginId", "read pluginid failed")
		return
	}

	// PluginRecord
	record := new(plugin.PluginRecord)
	err := ctx.ReadJSON(record)
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

func ModifyPluginRecord(ctx context.Context) {
	uid := ctx.Values().GetString(CTX_USERID)
	if uid == "" {
		sendNoUserResponse(ctx)
		return
	}

	pid := ctx.Params().Get("pid")
	if pid == "" {
		SendResponse(ctx, http.StatusInternalServerError, "unknownPluginId", "read pluginid failed")
		return
	}

	rid := ctx.Params().Get("rid")
	if rid == "" {
		SendResponse(ctx, http.StatusInternalServerError, "unknownRecordId", "read recordid failed")
		return
	}

	// PluginRecord
	pmap := make(map[string]interface{})
	err := ctx.ReadJSON(pmap)
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

func DeletePluginRecord(ctx context.Context) {
	uid := ctx.Values().GetString(CTX_USERID)
	if uid == "" {
		return
	}

	pid := ctx.Params().Get("pid")
	if pid == "" {
		SendResponse(ctx, http.StatusInternalServerError, "unknownPluginId", "read pluginid failed")
		return
	}

	rid := ctx.Params().Get("rid")
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

func sendNoUserResponse(ctx context.Context) {
	SendResponse(ctx, http.StatusInternalServerError, "unknownUser", "get userid from context failed")
}
