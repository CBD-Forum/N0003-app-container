# Introduction
Currently, the users could be classified into four kinds, 
  - the client who is responsible for creating a new shipping order,
  - the cargo agent who should book space from shippers and book vehicles from carriers after receiving a shipping order delegated by the client.
  - the shipper who owns the vessels and is in charge of managing containers and making shipping schedules.
  - the carrier who possesses the vehicles which are used to carrying goods in land transport.

The chaincode mainly functions in three ways,
  - by simulating the data flow of the shipping orders,
  - by informing of messages the users when order failed or order processed,
  - by managing resources like containers and vehicles.

