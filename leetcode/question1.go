/*
Tesla ships over-the-air software updates to vehicles every few weeks.
A good portion of this data is sent over cell networks, which is costly.
Imagine a world where Tesla has implemented a peer-to-peer data transfer system
that enables vehicles to share data with other vehicles over wifi in unique 10MB chunks.
Server-side, we would have a representation of a vehicle that would receive a callback whenever
a vehicle completed sending a chunk to another vehicle.
Assume this system is now built and running.
We have a big release coming up. You need to minimize cell costs by building a method rankVehicles
that returns a list of all vehicles in the network prioritized by which are most important to
"seed" with data via cell downloads.
For now we can assume importance is measured by a vehicle's historical contribution to the
peer-to-peer network: all direct sends + all indirect sends
Example:
Assume vehicles A, B, C, D
Assume chunks 1, 2, 3, 4
Assume Server sends chunks [1,2] to A & chunks [3,4] to C
A sends [1, 2] to B
A sends [1] to C
C sends [4] to D
C sends [3] to B
B sends [3] to D
D sends [3] to A
A's total contribution should be 30 (MB)
B's total contribution should be 20
C's total contribution should be 40
D's total contribution should be 10
rankVehicles should return [C,A,B,D]
*/

package main

import (
	"fmt"
	"log"
	"reflect"
)

var server *VehicleServer

func main() {
	server = NewVehicleServer()

	vA := server.StartUpdate("A", []int{1, 2})
	vB := server.StartUpdate("B", []int{})
	vC := server.StartUpdate("C", []int{3, 4})
	vD := server.StartUpdate("D", []int{})

	vA.CompleteChunkTransaction("B", 1)
	vA.CompleteChunkTransaction("B", 2)
	vA.CompleteChunkTransaction("C", 1)

	vC.CompleteChunkTransaction("D", 4)
	vC.CompleteChunkTransaction("B", 3)

	vB.CompleteChunkTransaction("D", 3)

	vD.CompleteChunkTransaction("A", 3)

	validate(server)
}

func NewVehicle(id string, chunks []int) *Vehicle {
	return &Vehicle{ID: id, chunks: chunks}
}

type Vehicle struct {
	ID     string
	chunks []int
}

// TODO: define this method
func (v *Vehicle) CompleteChunkTransaction(vehicleID string, chunkID int) {
	server.RecordChunkTransaction(v.ID, vehicleID, chunkID)
}

func NewVehicleServer() *VehicleServer {
	return &VehicleServer{}
}

type VehicleServer struct {
	// TODO: define proper data structure
}

func validate(server *VehicleServer) {
	if !reflect.DeepEqual(server.RankVehicles(), []string{"C", "A", "B", "D"}) {
		log.Fatal("Inaccurate ordering:", server.RankVehicles())
	}

	fmt.Println("Success!")
}

func (s *VehicleServer) StartUpdate(vehicleID string, chunkIDS []int) *Vehicle {
	return NewVehicle(vehicleID, chunkIDS)
}

func (s *VehicleServer) RecordChunkTransaction(sender, recipient string, chunkID int) {
	// TODO: implement record_chunk_transaction business logic to allow proper rank_vehicles computation
}

func (s *VehicleServer) RankVehicles() (vehicleIDS []string) {
	// TODO: implement rank_vehicles, should return a list of vehicle ids ranked by their total contribution

	return vehicleIDS
}
