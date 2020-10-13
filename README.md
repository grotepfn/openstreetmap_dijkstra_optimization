# Dijkstra Water Routing Optimization

## Usage
Create a bitarray for routable points on water with the help of coastlines.go.
Provide an OSM file with valid closed polygons.
The bitarray then is saved as file bitArray in folder data.
```bash
go run cmd\coastlines\main.go
```

Create routable shortcut squares within the bitarray with help of main.go.
```bash
go run cmd\optimization\main.go
```
This go program produces a three dimensional array with the optimization squares and 
saves it to the file optimization_squares.
The squares are stored in the form of [[[y,x],[y2,x2]]] where the upper left point and the lower right point of the squares are stored. The value of y can be seen as the vertical distance to the upper left point (0,0) of the bitarray, x represents the corresponding horizontal distance.
Test this file with help of main_test.go
```bash
go test cmd\optimization\main.go cmd\optimization\main_test.go
```
Also, the file optimization_squares_distances is produced. This file contains a three dimensional array with the same length of optimization_squares. The distances from each border point to each border point of the squares are stored, first all points of the upper x-axis and lower x-axis border. Second the points of the left y-axis and right y-axis borderline.

Start the golang routing backend with
```bash
go run main.go
```

Open the index.html with a webbrowser to use the application.

Run the benchmark comparing different routing algorithms with 
```bash
go run cmd\benchmark\main.go
```

A visualized example with help of geojson.io of outcome optimization squares can be found below.
![Optimazion squares visualized with geojson.io](https://github.com/grotepfn/openstreetmap_dijkstra_optimization/blob/master/data/dc1f91a2f0a5412b0413bbac4324d057.png?raw=true)
