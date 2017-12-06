package graph

import (
	"errors"
	//	"fmt"
	"math"
)

type Vertex struct {
	NextVertex *Vertex
	Key        string
	Arc        *Arc
	InDegree   int
	OutDegree  int

	Processed  bool    // for BrFS, DFS
	Parent     *Vertex // to trace when a solution is a path, not a state
	InTree     bool
	PathLength float64
}

type Arc struct {
	Dest    *Vertex
	NextArc *Arc
	Weight  float64
	InTree  bool
}

func NewVertex() *Vertex {
	return &Vertex{NextVertex: nil, Arc: nil, Key: "", InDegree: 0, OutDegree: 0, Processed: false, Parent: nil}
}

func NewArc(weight float64) *Arc {
	return &Arc{Dest: nil, NextArc: nil, Weight: weight}
}

type Graph struct {
	First *Vertex
	Count int
}

func NewGraph() *Graph {
	return &Graph{First: nil, Count: 0}
}

func (graph *Graph) InsertVertex(dataKey string) {
	newPtr := NewVertex()
	newPtr.Key = dataKey
	graph.Count++
	if graph.First == nil {
		graph.First = newPtr
	} else {
		locPtr := graph.First
		var prePtr *Vertex = nil
		for locPtr != nil && dataKey > locPtr.Key {
			prePtr = locPtr
			locPtr = locPtr.NextVertex
		}
		if prePtr == nil {
			graph.First = newPtr
		} else {
			prePtr.NextVertex = newPtr
		}
		newPtr.NextVertex = locPtr
	}
}

func (graph *Graph) DeleteVertex(dataKey string) {
	var prePtr, locPtr *Vertex
	if graph != nil {
		prePtr = nil
		locPtr = graph.First
		for locPtr != nil && dataKey > locPtr.Key {
			prePtr = locPtr
			locPtr = locPtr.NextVertex
		}
		if locPtr == nil || dataKey != locPtr.Key {
			return // not found
		}
		if locPtr.InDegree > 0 || locPtr.OutDegree > 0 {
			return // delete only when degree is 0
		}
		if prePtr == nil { // first vertex will be deleted
			graph.First = locPtr.NextVertex
		} else {
			prePtr.NextVertex = locPtr.NextVertex
		}
		graph.Count--
	}
}

func (graph *Graph) InsertArc(fromKey, toKey string, weight float64) error {
	var fromPtr *Vertex = graph.First
	for fromPtr != nil && fromKey > fromPtr.Key {
		fromPtr = fromPtr.NextVertex
	}
	if fromPtr == nil || fromKey != fromPtr.Key {
		return errors.New("FromKey not found")
	}
	var toPtr *Vertex = graph.First
	for toPtr != nil && toKey > toPtr.Key {
		toPtr = toPtr.NextVertex
	}
	if toPtr == nil || toKey != toPtr.Key {
		return errors.New("ToKey not found")
	}
	var newArc *Arc = NewArc(weight)
	newArc.Dest = toPtr
	fromPtr.OutDegree++
	toPtr.InDegree++
	if fromPtr.Arc == nil {
		fromPtr.Arc = newArc
		newArc.NextArc = nil
	} else {
		var arcPrePtr *Arc = nil
		var arcWalkPtr *Arc = fromPtr.Arc
		for arcWalkPtr != nil && toKey > arcWalkPtr.Dest.Key {
			arcPrePtr = arcWalkPtr
			arcWalkPtr = arcWalkPtr.NextArc
		}
		if arcPrePtr == nil {
			fromPtr.Arc = newArc
		} else {
			arcPrePtr.NextArc = newArc
		}
		newArc.NextArc = arcWalkPtr
	}
	return nil
}

type QueueNode struct {
	Data *Vertex
	Next *QueueNode
}

func NewQueueNode() *QueueNode {
	return &QueueNode{Data: nil, Next: nil}
}

type Queue struct {
	Front           *QueueNode
	Rear            *QueueNode
	Count           int
	IsPriorityQueue bool
}

func NewQueue(isPriorityQueue bool) *Queue {
	return &Queue{Front: nil, Rear: nil, Count: 0, IsPriorityQueue: isPriorityQueue}
}

func (queue *Queue) Enqueue(data *Vertex) {
	var newPtr *QueueNode = NewQueueNode()
	newPtr.Data = data
	newPtr.Next = nil
	if queue.Count == 0 {
		queue.Front = newPtr
		queue.Rear = newPtr
	} else {
		if !queue.IsPriorityQueue {
			queue.Rear.Next = newPtr
			queue.Rear = newPtr
		} else {
			if data == nil { // do not accept nil data in priority queue
				return
			}
			ptr := queue.Front
			// search to find a correct place to put newPtr
			var prePtr *QueueNode
			prePtr = nil
			for ptr != nil && data.PathLength > ptr.Data.PathLength {
				prePtr = ptr
				ptr = ptr.Next
			}
			if ptr == nil {
				queue.Rear = newPtr
			}
			if prePtr == nil {
				queue.Front = newPtr
			} else {
				prePtr.Next = newPtr
			}
			newPtr.Next = ptr
		}
	}
	queue.Count++
}

func (queue *Queue) Dequeue(dataOut **Vertex) error {
	if queue.Count == 0 {
		return errors.New("Queue is empty")
	}
	*dataOut = queue.Front.Data
	if queue.Count == 1 {
		queue.Rear = nil
	}
	queue.Front = queue.Front.Next
	queue.Count--
	return nil
}

func (queue *Queue) IsEmpty() bool {
	return queue.Count == 0
}

func (queue *Queue) GetFront() *Vertex {
	if queue.Count == 0 {
		return nil
	}
	return queue.Front.Data
}

func (queue *Queue) GetRear() *Vertex {
	if queue.Count == 0 {
		return nil
	}
	return queue.Rear.Data
}

/* STACK */

type StackNode struct {
	Data *Vertex
	Next *StackNode
}

type Stack struct {
	Top   *StackNode
	Count int
}

func NewStack() *Stack {
	return &Stack{Top: nil, Count: 0}
}

func (stack *Stack) Push(dataIn *Vertex) {
	pNew := &StackNode{Data: dataIn, Next: nil}
	pNew.Next = stack.Top
	stack.Top = pNew
	stack.Count++
}

func (stack *Stack) Pop() *Vertex {
	if stack.Count == 0 {
		return nil
	}
	//deletePtr = stack.Top // Keep to recycle the memory. Not neccessary in Go
	dataOut := stack.Top.Data
	stack.Top = stack.Top.Next
	stack.Count--
	return dataOut
}

func (stack *Stack) IsEmpty() bool {
	return stack.Count == 0
}

func (graph *Graph) FindDistance(route []string) (float64, error) {
	var vertexRoute []*Vertex = make([]*Vertex, len(route))
	// get vertex pointers for the keys
	for index, vertexKey := range route {
		ptr := graph.First
		for ptr != nil {
			if ptr.Key == vertexKey {
				vertexRoute[index] = ptr
				break
			}
			ptr = ptr.NextVertex
		}
		if ptr == nil {
			return 0, errors.New("Vertex Key not found")
		}
	}
	// get weights and calculate results
	distance := 0.0
	for index, vertexPtr := range vertexRoute {
		if index < len(vertexRoute)-1 {
			arcPtr := vertexPtr.Arc
			for arcPtr != nil {
				if arcPtr.Dest.Key == vertexRoute[index+1].Key {
					distance += arcPtr.Weight
					break
				}
				arcPtr = arcPtr.NextArc
			}
			if arcPtr == nil {
				return 0, errors.New("NO SUCH ROUTE")
			}
		}
	}
	return distance, nil
}

// Find round trips from the vertex 'vertexKey' with max number of stops 'stops'
func (graph *Graph) FindRoundTripWithMaxStops(vertexKey string, stops int) ([]string, error) {
	// Find vertex with the key
	if graph.First == nil {
		return nil, errors.New("Graph is empty") // graph is empty
	}
	vPtr := graph.First
	for vPtr != nil && vPtr.Key != vertexKey {
		vPtr = vPtr.NextVertex
	}
	if vPtr == nil {
		return nil, errors.New("Key not found") // vertex Key not found
	}
	// start to do a breadth-first search
	vPtr.Parent = nil
	queue := NewQueue(false)
	queue.Enqueue(vPtr)
	queue.Enqueue(nil)
	level := 0
	for !queue.IsEmpty() {
		queue.Dequeue(&vPtr)
		if vPtr == nil { // this is a separator to separate the vertexes with the same number of stops in queue
			level++
			if !queue.IsEmpty() {
				queue.Enqueue(nil) // insert nil as a separator
			}
		} else {
			if level > 0 && vPtr.Key == vertexKey { // found a solution
				//fmt.Println("Found a solution at level", level)
				return processSolution(vPtr), nil
			} else {
				if level >= stops { // if level is greater than stops, do not continue to consider next vertexes
					continue
				}
				aPtr := vPtr.Arc
				for aPtr != nil {
					// get next vertex and put it into queue
					data := aPtr.Dest

					dataPtr := &Vertex{Key: data.Key,
						NextVertex: data.NextVertex,
						Arc:        data.Arc,
						InDegree:   data.InDegree,
						OutDegree:  data.OutDegree,
						Processed:  data.Processed,
						Parent:     data.Parent}
					dataPtr.Parent = vPtr // update the parent pointer to trace to the root vertex in the solution if there is a path
					queue.Enqueue(dataPtr)
					//fmt.Println("Enqueue from ", vPtr.Data, " dest ", dataPtr.Data)
					aPtr = aPtr.NextArc
				}
			}
		}
	}
	return nil, errors.New("Result not found")
}

// Process a solution by tracing the parent pointer to print a path on screen.
func processSolution(vertex *Vertex) []string {
	stack := NewStack()

	ptr := vertex
	for ptr != nil {
		stack.Push(ptr)
		//fmt.Println(ptr.Data, "-")
		ptr = ptr.Parent
	}
	var result []string
	result = make([]string, stack.Count)
	i := 0
	for !stack.IsEmpty() {
		ptr := stack.Pop()
		result[i] = ptr.Key
		i++
		//if stack.Count > 0 {
		//	fmt.Print("-")
		//}
	}
	//fmt.Println()
	return result
}

// Find trips from 'fromKey' to 'toKey' with a given number of stops 'stops'
func (graph *Graph) FindTripExactStops(fromKey, toKey string, stops int) ([]string, error) {
	// Find vertex with the key
	if graph.First == nil {
		return nil, errors.New("Graph is empty")
	}
	vFromPtr := graph.First
	for vFromPtr != nil && vFromPtr.Key != fromKey {
		vFromPtr = vFromPtr.NextVertex
	}
	if vFromPtr == nil {
		return nil, errors.New("FromKey not found")
	}
	vToPtr := graph.First
	for vToPtr != nil && vToPtr.Key != toKey {
		vToPtr = vToPtr.NextVertex
	}
	if vToPtr == nil {
		return nil, errors.New("ToKey not found")
	}
	// use breadth-first search to traverse a graph by level
	vPtr := vFromPtr
	vPtr.Parent = nil
	queue := NewQueue(false)
	queue.Enqueue(vPtr)
	queue.Enqueue(nil)
	level := 0
	for !queue.IsEmpty() {
		queue.Dequeue(&vPtr)
		if vPtr == nil { // this is a separator
			level++
			if !queue.IsEmpty() {
				queue.Enqueue(nil)
			}
		} else {
			if level > 0 && level == stops && vPtr.Key == toKey { // found a solution
				//fmt.Println("Found a solution at level", level)
				return processSolution(vPtr), nil
			} else {
				if level >= stops {
					continue
				}
				aPtr := vPtr.Arc
				for aPtr != nil {
					data := aPtr.Dest

					dataPtr := &Vertex{
						Key:        data.Key,
						NextVertex: data.NextVertex,
						Arc:        data.Arc,
						InDegree:   data.InDegree,
						OutDegree:  data.OutDegree,
						Processed:  data.Processed,
						Parent:     data.Parent}
					dataPtr.Parent = vPtr // update parent to trace the solution if it is a path
					queue.Enqueue(dataPtr)
					//fmt.Println("Enqueue from ", vPtr.Data, " dest ", dataPtr.Data)
					aPtr = aPtr.NextArc
				}
			}
		}
	}
	return nil, errors.New("No solution found")
}

// This method uses Best-First-Search algorithm with the help of a priority queue.
func (graph *Graph) FindShortestRoute(fromKey string, toKey string) ([]string, error) {
	// Find vertex with the key
	if graph.First == nil {
		return nil, errors.New("Graph is empty")
	}
	vFromPtr := graph.First
	for vFromPtr != nil && vFromPtr.Key != fromKey {
		vFromPtr = vFromPtr.NextVertex
	}
	if vFromPtr == nil {
		return nil, errors.New("FromKey not found")
	}
	vToPtr := graph.First
	for vToPtr != nil && vToPtr.Key != toKey {
		vToPtr = vToPtr.NextVertex
	}
	if vToPtr == nil {
		return nil, errors.New("ToKey not found")
	}

	vPtr := graph.First
	for vPtr != nil {
		vPtr.PathLength = math.MaxFloat64 // PathLength keeps distance from the source to the current vertex. It is initialised to INFINITY
		vPtr.Processed = false
		vPtr = vPtr.NextVertex
	}

	// use Uniform-cost search to traverse a graph
	vPtr = vFromPtr
	vPtr.PathLength = 0.0 // distance from source to source is 0

	vPtr.Parent = nil
	queue := NewQueue(true)
	queue.Enqueue(vPtr)
	vPtr.Processed = true // true means that it is used to be in the queue
	for !queue.IsEmpty() {
		queue.Dequeue(&vPtr)
		//fmt.Println("Dequeue:", vPtr.Key)
		if vPtr.Key == toKey {
			return processSolution(vPtr), nil
		}
		aPtr := vPtr.Arc
		for aPtr != nil {
			dest := aPtr.Dest
			if !dest.Processed {
				dest.PathLength = vPtr.PathLength + aPtr.Weight
				dest.Parent = vPtr // update parent to trace the solution
				queue.Enqueue(dest)
				dest.Processed = true
			} else {
				if dest.PathLength > vPtr.PathLength+aPtr.Weight {
					dest.PathLength = vPtr.PathLength + aPtr.Weight
					dest.Parent = vPtr
					queue.Enqueue(dest)
					dest.Processed = true
				}
			}
			aPtr = aPtr.NextArc
		}
	}
	return nil, errors.New("No solution found")
}

// This method uses Best-First-Search algorithm with the help of a priority queue.
func (graph *Graph) FindShortestRoundTrip(fromKey string) ([]string, error) {
	// Find vertex with the key
	if graph.First == nil {
		return nil, errors.New("Graph is empty")
	}
	vFromPtr := graph.First
	for vFromPtr != nil && vFromPtr.Key != fromKey {
		vFromPtr = vFromPtr.NextVertex
	}
	if vFromPtr == nil {
		return nil, errors.New("FromKey not found")
	}

	vPtr := graph.First
	for vPtr != nil {
		vPtr.PathLength = math.MaxFloat64 // PathLength keeps distance from the source to the current vertex. It is initialised to INFINITY
		vPtr.Processed = false
		vPtr = vPtr.NextVertex
	}

	// use Uniform-cost search to traverse a graph
	vPtr = vFromPtr
	vPtr.PathLength = 0.0 // distance from source to source is 0

	queue := NewQueue(true)
	queue.Enqueue(vPtr)
	vPtr.Processed = true // true means that it is used to be in the queue
	for !queue.IsEmpty() {
		queue.Dequeue(&vPtr)
		//fmt.Println("Dequeue:", vPtr.Key, "PathLength:", vPtr.PathLength)
		if vPtr.PathLength > 0 && vPtr.Key == fromKey {
			// process the solution
			ptr := vPtr.Parent
			vPtr.Parent = nil
			vPtr = ptr
			solution := append(processSolution(vPtr), fromKey)
			return solution, nil
		}
		aPtr := vPtr.Arc
		for aPtr != nil {
			dest := aPtr.Dest
			if !dest.Processed || dest.PathLength == 0 {
				dest.PathLength = vPtr.PathLength + aPtr.Weight
				dest.Parent = vPtr // update parent to trace the solution
				queue.Enqueue(dest)
				dest.Processed = true
				//fmt.Println("Parent of", dest.Key, "is", vPtr.Key)
			} else {
				if dest.PathLength > vPtr.PathLength+aPtr.Weight {
					dest.PathLength = vPtr.PathLength + aPtr.Weight
					dest.Parent = vPtr
					queue.Enqueue(dest)
					dest.Processed = true
					//fmt.Println("Parent of", dest.Key, "is", vPtr.Key)
				}
			}
			aPtr = aPtr.NextArc
		}
	}
	return nil, errors.New("No solution found")
}

func (graph *Graph) FindRoundTripWithMaxWeight(vertexKey string, maxWeight float64) ([][]string, error) {
	// Find vertex with the key
	var solutions [][]string

	if graph.First == nil {
		return nil, errors.New("Graph is empty")
	}
	vPtr := graph.First
	for vPtr != nil && vPtr.Key != vertexKey {
		vPtr = vPtr.NextVertex
	}
	if vPtr == nil {
		return nil, errors.New("Vertex Key not found")
	}
	vPtr.Parent = nil
	vPtr.PathLength = 0
	queue := NewQueue(false)
	queue.Enqueue(vPtr)
	for !queue.IsEmpty() {
		queue.Dequeue(&vPtr)
		if vPtr.PathLength < maxWeight {
			if vPtr.PathLength > 0 && vPtr.Key == vertexKey { // found a solution
				//fmt.Println("Found a solution")
				solution := processSolution(vPtr)
				solutions = append(solutions, solution)
			}
			aPtr := vPtr.Arc
			for aPtr != nil {
				dest := aPtr.Dest

				dataPtr := &Vertex{Key: dest.Key, // Need to copy because we might have several solutions
					NextVertex: dest.NextVertex,
					Arc:        dest.Arc,
					InDegree:   dest.InDegree,
					OutDegree:  dest.OutDegree,
					Processed:  dest.Processed,
					Parent:     dest.Parent,
					PathLength: vPtr.PathLength + aPtr.Weight}
				dataPtr.Parent = vPtr // update parent to trace the solution if it is a path
				if dataPtr.PathLength < maxWeight {
					queue.Enqueue(dataPtr)
				}
				//fmt.Println("Enqueue from ", vPtr.Data, " dest ", dataPtr.Data)
				aPtr = aPtr.NextArc
			}

		}
	}
	return solutions, nil
}
