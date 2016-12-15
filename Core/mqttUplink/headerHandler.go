package mqttUplink

import (
	"encoding/binary"
	"errors"

	components "github.com/LoRaWanSoFa/LoRaWanSoFa/Components"
	"github.com/LoRaWanSoFa/LoRaWanSoFa/DBC/DatabaseConnector"
)

type HeaderHandler interface {
	CreateNewHeader(payload []byte, devEUI string) ([]components.Sensor, error)
	StoreHeader(header []components.Sensor, devEUI string) error
}

type headerHandler struct {
}

func NewHeaderHandler() HeaderHandler {
	h := new(headerHandler)
	return h
}

func (h *headerHandler) CreateNewHeader(payload []byte, devEUI string) ([]components.Sensor, error) {
	var sensors []components.Sensor
	if h.checkLength(payload) {
		for i := 1; i < len(payload); i = i + 3 {
			sensor := h.createSensor(payload[i : i+3])
			sensor.HeaderOrder = (i + 2) / 3
			sensors = append(sensors, sensor)
		}
	} else {
		err := errors.New("Header of unkown length was send.")
		return nil, err
	}
	return sensors, nil
}

func (h *headerHandler) StoreHeader(newHeader []components.Sensor, devEUI string) error {
	oldHeader, err := DatabaseConnector.GetFullHeader(devEUI)
	if err != nil {
		return err
	}
	var headerOrderChanges, activationChanged []components.Sensor

	for i := range oldHeader {
		if b, position := h.containsSensor(oldHeader[i], newHeader); b {
			sensor := newHeader[position]
			newHeader = append(newHeader[:i], newHeader[i+1:]...)
			oldHeader[i].Soft_deleted = false
			activationChanged = append(activationChanged, oldHeader[i])
			if sensor.HeaderOrder != oldHeader[i].HeaderOrder {
				oldHeader[i].HeaderOrder = sensor.HeaderOrder
				headerOrderChanges = append(headerOrderChanges, oldHeader[i])
			}
		} else if !oldHeader[i].Soft_deleted {
			oldHeader[i].Soft_deleted = true
			activationChanged = append(activationChanged, oldHeader[i])
		}
	}
	if len(newHeader) != 0 {
		//TODO: DatabaseConnector.NewSensor(newHeader)
	}
	if len(activationChanged) != 0 {
		//TODO: DatabaseConnector.ChangeSensorActivationStates(activationChanged)
	}
	if len(headerOrderChanges) != 0 {
		//TODO: DatabaseConnector.ChangeSensorOrder(headerOrderChanges)
	}

	return nil
}

func (h *headerHandler) containsSensor(sensor components.Sensor, sensors []components.Sensor) (bool, int) {
	for i := range sensors {
		if sensor.SameSensor(sensors[i]) {
			return true, i
		}
	}
	return false, 0
}

func (h *headerHandler) createSensor(payload []byte) components.Sensor {
	var sensor components.Sensor

	ioType := payload[0]
	ioType = ioType >> 6
	sensor.IoType = int(ioType)

	ioAddress := payload[0]
	ioAddress = ioAddress << 2 >> 4
	sensor.IoAddress = int(ioAddress)

	sensorType := payload[:2]
	sensorType[0] = sensorType[0] << 6 >> 6
	sensor.SensorType = int(binary.BigEndian.Uint16(sensorType))

	lenghtOfValues := payload[2]
	lenghtOfValues = lenghtOfValues >> 5
	if lenghtOfValues == 0 {
		lenghtOfValues = 1
	}
	sensor.LenghtOfValues = int(lenghtOfValues)

	numberOfValues := payload[2]
	numberOfValues = numberOfValues << 2 >> 5
	sensor.NumberOfValues = int(numberOfValues)

	sensor.Soft_deleted = false

	return sensor
}

func (hc *headerHandler) checkLength(payload []byte) bool {
	return len(payload)%3 == 1
}
