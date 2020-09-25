# Dijkstra Water Routing Optimization

## Usage
Create a bitarray for routable points on water with the help of coastlines.go.
Provide an OSM file with valid closed polygons.
The bitarray then is saved as file bitArray in folder data.
```bash
go run coastlines.go
```

Create routable shortcut squares within the bitarray with help of optimization.go.
```bash
go run optimization.go
```
This go program produces a three dimensional array with the optimization squares and 
saves it to the file optimization_squares.
The squares are stored in the form of [[[y,x],[y2,x2]]] whereas the upper left point and the lower right point of the squares are stored. The value of y can be seen as the vertical distance to the upper left point (0,0) of the bitarray, x represents the corresponding horizontal distance.
Test this file with help of optimization_test.go
```bash
go test optimization.go optimization_test.go
```