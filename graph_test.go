package graph_test

import (
	"fmt"
	"maxsure.org/demo/graph"
	"reflect"
	"testing"
)

func initGraph(graph *graph.Graph) {
	graph.InsertVertex("a")
	graph.InsertVertex("b")
	graph.InsertVertex("c")
	graph.InsertVertex("d")
	graph.InsertVertex("e")
	graph.InsertArc("a", "b", 5)
	graph.InsertArc("b", "c", 4)
	graph.InsertArc("c", "d", 8)
	graph.InsertArc("d", "c", 8)
	graph.InsertArc("d", "e", 6)
	graph.InsertArc("a", "d", 5)
	graph.InsertArc("c", "e", 2)
	graph.InsertArc("e", "b", 3)
	graph.InsertArc("a", "e", 7)
}

func TestProblem1(t *testing.T) {
	graph := graph.NewGraph()
	initGraph(graph)
	fmt.Println("Testing problem 1: Find distance of a route")
	result, err := graph.FindDistance([]string{"a", "b", "c"})
	if err != nil {
		t.Errorf("No such route")
	}
	expected := 9.0
	if result != expected {
		t.Errorf("The result should be 9")
	}
}

func TestProblem2(t *testing.T) {
	fmt.Println("Testing problem 2: Find distance of a route")
	graph := graph.NewGraph()
	initGraph(graph)

	result, err := graph.FindDistance([]string{"a", "d"})
	if err != nil {
		t.Errorf("No such route")
	}
	expected := 5.0
	if result != expected {
		t.Errorf("The result should be 5")
	}
}

func TestProblem3(t *testing.T) {
	fmt.Println("Testing problem 3: Find distance of a route")
	graph := graph.NewGraph()
	initGraph(graph)

	result, err := graph.FindDistance([]string{"a", "d", "c"})
	if err != nil {
		t.Errorf("No such route")
	}
	expected := 13.0
	if result != expected {
		t.Errorf("The result should be 13")
	}
}

func TestProblem4(t *testing.T) {
	fmt.Println("Testing problem 4: Find distance of a route")
	graph := graph.NewGraph()
	initGraph(graph)

	result, err := graph.FindDistance([]string{"a", "e", "b", "c", "d"})
	if err != nil {
		t.Errorf("No such route")
	}
	expected := 22.0
	if result != expected {
		t.Errorf("The result should be 22")
	}
}

func TestProblem5(t *testing.T) {
	fmt.Println("Testing problem 5: Find distance of a route")
	graph := graph.NewGraph()
	initGraph(graph)

	_, err := graph.FindDistance([]string{"a", "e", "d"})
	if err == nil {
		t.Errorf("Error should NOT be NIL")
	}
	expected := "NO SUCH ROUTE"
	if err.Error() != expected {
		t.Errorf("Error message should not be 'No such route' ", err.Error())
	}
}

func TestProblem6(t *testing.T) {
	fmt.Println("Testing problem 6: Find round trip with max stops")
	from := "c"
	stops := 3
	fmt.Println("From", from, "Stops", stops)
	graph := graph.NewGraph()
	initGraph(graph)

	result, err := graph.FindRoundTripWithMaxStops(from, stops)
	fmt.Println(result)
	if err != nil {
		t.Errorf("Error should be NIL")
	}
	if result == nil {
		t.Errorf("Error: No solution")
	}
}

func TestProblem7(t *testing.T) {
	fmt.Println("Testing problem 7: Find a trip with a number of stops")
	from := "a"
	to := "c"
	stops := 4
	fmt.Println("From", from, "To", to, "Stops", stops)
	graph := graph.NewGraph()
	initGraph(graph)

	result, err := graph.FindTripExactStops(from, to, stops)
	fmt.Println(result)
	if err != nil {
		t.Errorf("Error should be NIL")
	}
	if result == nil {
		t.Errorf("Error: No solution")
	}
}

func TestProblem8(t *testing.T) {
	fmt.Println("Testing problem 8: Find the shortest route")
	from := "a"
	to := "c"
	fmt.Println("From", from, "To", to)
	graph := graph.NewGraph()
	initGraph(graph)

	result, err := graph.FindShortestRoute(from, to)
	fmt.Println(result)
	if err != nil {
		t.Errorf("Error should be NIL")
	}
	if result == nil {
		t.Errorf("Error: No solution")
	}
	expected := []string{"a", "b", "c"}
	eq := reflect.DeepEqual(result, expected)
	if !eq {
		t.Errorf("Solution is not correct")
	}
}

func TestProblem9(t *testing.T) {
	fmt.Println("Testing problem 9: Find the shortest round trip")
	from := "c"
	fmt.Println("From", from)
	graph := graph.NewGraph()
	initGraph(graph)

	result, err := graph.FindShortestRoundTrip(from)
	fmt.Println(result)
	if err != nil {
		t.Errorf("Error should be NIL")
	}
	if result == nil {
		t.Errorf("Error: No solution")
	}
}

func TestProblem10(t *testing.T) {
	fmt.Println("Testing problem 10: Find round trips with a max weight")
	from := "c"
	weight := 30.0
	fmt.Println("From", from, "Weight", weight)
	graph := graph.NewGraph()
	initGraph(graph)

	result, err := graph.FindRoundTripWithMaxWeight(from, weight)
	fmt.Println(result)
	if err != nil {
		t.Errorf("Error should be NIL")
	}
	expected := [][]string{
		[]string{"c", "d", "c"},
		[]string{"c", "e", "b", "c"},
		[]string{"c", "d", "e", "b", "c"},
		[]string{"c", "d", "c", "e", "b", "c"},
		[]string{"c", "e", "b", "c", "d", "c"},
		[]string{"c", "e", "b", "c", "e", "b", "c"},
		[]string{"c", "e", "b", "c", "e", "b", "c", "e", "b", "c"}}
	eq := reflect.DeepEqual(result, expected)
	if !eq {
		fmt.Println(result)
		t.Errorf("Solution is not correct")
	}
}

func TestPriorityQueue(t *testing.T) {
	fmt.Println("Testing priority queue")
	queue := graph.NewQueue(true)

	vPtr := &graph.Vertex{Key: "a", PathLength: 8.0}
	queue.Enqueue(vPtr)
	vPtr = &graph.Vertex{Key: "b", PathLength: 5.0}
	queue.Enqueue(vPtr)
	vPtr = &graph.Vertex{Key: "c", PathLength: 3.0}
	queue.Enqueue(vPtr)
	vPtr = &graph.Vertex{Key: "d", PathLength: 6.0}
	queue.Enqueue(vPtr)

	queue.Dequeue(&vPtr)
	if vPtr.PathLength != 3 {
		t.Errorf("Error in priority queue")
	}
	queue.Dequeue(&vPtr)
	if vPtr.PathLength != 5 {
		t.Errorf("Error in priority queue")
	}
}
