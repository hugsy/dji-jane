package djijane

import (
	"dji-joe"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func ExtractJsonResponseAs(r *http.Request, t string) (interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		Log.ErrorF("Invalid body: %+v", err)
		return nil, err
	}
	r.Body.Close()

	var msg interface{}

	switch t {
	case "heartbeat":
		msg = new(djijoe.HeartBeatMessage)

	case "wakeup":
		msg = new(djijoe.WakeUpMessage)

	case "shutdown":
		msg = new(djijoe.ShutdownMessage)

	case "info":
		msg = new(djijoe.DroneInfoMessage)

	default:
		Log.FatalF("incorrect type: %s", t)
	}

	err = json.Unmarshal(body, msg)
	if err != nil {
		Log.ErrorF("Failed to parse JSON: %+v", err)
		return nil, err
	}
	return msg, nil
}

/*
Probe heartbeat handler

test with
curl -i -d '{"ts":"0001-01-01T00:00:00Z", "host":"foo"}' -X POST http://localhost:8000/api/heartbeat
*/
func Heartbeat(w http.ResponseWriter, r *http.Request) {
	msg, err := ExtractJsonResponseAs(r, "heartbeat")
	if err != nil {
		return
	}
	heartbeat := msg.(*djijoe.HeartBeatMessage)

	Log.NoticeF("Received HEARTBEAT from '%s'", heartbeat.Hostname)

	i, p := djijoe.GetProbeByName(probes, heartbeat.Hostname)
	if i == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	p.LastHeartbeat = time.Now()
	w.WriteHeader(http.StatusAccepted)
}

/*
Probe wake-up handler

test with
curl -i --data '{"ts":"0001-01-01T00:00:00Z","host":"foo"}' -X POST http://localhost:8000/api/wakeup
*/
func WakeUp(w http.ResponseWriter, r *http.Request) {
	msg, err := ExtractJsonResponseAs(r, "wakeup")
	if err != nil {
		return
	}
	wakeup := msg.(*djijoe.WakeUpMessage)

	idx, _ := djijoe.GetProbeByName(probes, wakeup.Hostname)
	if idx != -1 {
		Log.ErrorF("A probe named '%s' already exists", wakeup.Hostname)
		w.WriteHeader(http.StatusConflict)
		return
	}

	p := djijoe.Probe{
		Hostname:       wakeup.Hostname,
		StartTime:      wakeup.Timestamp,
		GpsCoordinates: wakeup.Position,
		LastHeartbeat:  time.Now(),
	}

	Log.NoticeF("Received WAKEUP from '%s' at %s (lat=%.3f long=%.3f)",
		p.Hostname, p.StartTime, p.GpsCoordinates.Lat(), p.GpsCoordinates.Lng())

	Log.DebugF("Added '%s' to `probes[]`", p)
	probes = append(probes, p)

	w.WriteHeader(http.StatusNoContent)
}

/*
Probe shutdown handler

test with
curl -i --data '{"ts":"0001-01-01T00:00:00Z","host":"foo","nb_beacon":1,"nb_probes":1}' -X POST http://localhost:8000/api/shutdown
*/
func ShutDown(w http.ResponseWriter, r *http.Request) {
	msg, err := ExtractJsonResponseAs(r, "shutdown")
	if err != nil {
		return
	}
	shutdown := msg.(*djijoe.ShutdownMessage)

	Log.NoticeF("Received SHUTDOWN from '%s' at %s: nb_beacons: %d, nb_probes: %d",
		shutdown.Hostname, shutdown.Timestamp,
		shutdown.BeaconFound, shutdown.ProbeRequestFound,
	)

	idx, cur_probe := djijoe.GetProbeByName(probes, shutdown.Hostname)
	if idx == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	Log.DebugF("Removed '%s' from `probes[]`", *cur_probe)
	probes, err = djijoe.RemoveProbe(probes, cur_probe)

	w.WriteHeader(http.StatusNoContent)
}

/*
New drone information handler

test with
curl -i --data '{"ts":"0001-01-01T00:00:00Z","host":"foo","type":1,"strength":-20,"vendor":"bar","macaddr":001122334455}' -X POST http://localhost:8000/api/info
*/
func NewDroneInfo(w http.ResponseWriter, r *http.Request) {
	msg, err := ExtractJsonResponseAs(r, "info")
	if err != nil {
		return
	}

	info := msg.(*djijoe.DroneInfoMessage)

	switch info.MessageType {
	case djijoe.TYPE_PROBE_REQUEST:
		Log.NoticeF("Received NEWDRONEINFO ProbeRequest by '%s' at %s from %s (%s) with strength: %d dBm",
			info.Hostname, info.Timestamp, info.MacAddress, info.Vendor, info.SignalStrength)

	case djijoe.TYPE_BEACON:
		Log.NoticeF("Received NEWDRONEINFO Beacon by '%s' at %s from %s (%s) with strength: %d dBm",
			info.Hostname, info.Timestamp, info.MacAddress, info.Vendor, info.SignalStrength)
	}

	w.WriteHeader(http.StatusAccepted)
}
