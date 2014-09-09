// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/golang/glog"
	//import the Paho Go MQTT library
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

//define a function for the default message handler
var f MQTT.MessageHandler = func(client *MQTT.MqttClient, msg MQTT.Message) {
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func client(id string) {
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions()
	// SetBroker("tcp://127.0.0.1:9900")
	opts.SetKeepAlive(10)
	opts.SetClientId("go-simple")
	opts.SetUsername(id)
	opts.SetPassword(psw)
	opts.SetStore(MQTT.NewMemoryStore())
	opts.AddBroker("tcp://127.0.0.1:9900")
	opts.SetDefaultPublishHandler(f)
	// opts.SetCleanSession(true)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	_, err := c.Start()
	if err != nil {
		glog.Info("Error：", err)
	}
}
