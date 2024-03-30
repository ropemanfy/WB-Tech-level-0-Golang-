package main

import "L0/nats-publisher/publisher"

func main() {
	pub := publisher.NewPublisher()
	pub.StartPublic(10)
}
