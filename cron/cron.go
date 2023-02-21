package cron

import (
	"log"
	"time"

	"github.com/potproject/ikuradon-salmon/dataaccess"
	"github.com/potproject/ikuradon-salmon/setting"
	"github.com/robfig/cron/v3"
)

func DeleteOldNotificationsCron() {
	p := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := p.Parse(setting.S.DeleteOldNotificationCron)
	if err != nil {
		log.Fatal(err)
		return
	}
	c := cron.New()
	c.AddFunc(setting.S.DeleteOldNotificationCron, deleteOldNotifications)
}

func deleteOldNotifications() {
	log.Println("Running DeleteOldNotificationsCron...")
	d, err := dataaccess.DA.ListAll()
	if err != nil {
		log.Fatal(err)
		return
	}
	count := 0
	now := time.Now().Unix()
	deleteDays := int64(setting.S.DeleteOldNotificationDays)
	for _, v := range d {
		lastUpDatedAt := v.Value.LastUpdatedAt
		if now-lastUpDatedAt < 60*60*24*deleteDays {
			continue
		}
		err = dataaccess.DA.Delete(v.Key)
		if err != nil {
			log.Print(err)
			return
		}
		count++
	}
	log.Printf("Done: Deleted %d notifications", count)
}
