ticker := time.NewTicker(5 * time.Second)
quit := make(chan struct{})
func() {
for {
select {
case err := <- ws.ReadJSON(RoomLastSecondsCheerCountMessage{}):
// do stuff
case <- quit:
ticker.Stop()
return
}
}
}()