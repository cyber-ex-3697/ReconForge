package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "html/template"
    "net/http"
    "os"
    "path/filepath"
)

// GraphData represents graph visualization data
type GraphData struct {
    Nodes []Node `json:"nodes"`
    Links []Link `json:"links"`
}

// Node represents a graph node
type Node struct {
    ID    string `json:"id"`
    Label string `json:"label"`
    Group string `json:"group"`
    Size  int    `json:"size"`
}

// Link represents a graph edge
type Link struct {
    Source string `json:"source"`
    Target string `json:"target"`
    Type   string `json:"type"`
}

func main() {
    var inputPath string
    var port int
    var openBrowser bool
    
    flag.StringVar(&inputPath, "input", "", "Input JSON file with graph data")
    flag.IntVar(&port, "port", 8080, "HTTP server port")
    flag.BoolVar(&openBrowser, "open", true, "Open browser automatically")
    flag.Parse()
    
    fmt.Println("=== ReconForge Graph Viewer ===\n")
    
    var graphData GraphData
    
    if inputPath != "" {
        // Load graph data from file
        data, err := os.ReadFile(inputPath)
        if err != nil {
            fmt.Printf("[!] Error reading file: %v\n", err)
            os.Exit(1)
        }
        
        if err := json.Unmarshal(data, &graphData); err != nil {
            fmt.Printf("[!] Error parsing JSON: %v\n", err)
            os.Exit(1)
        }
        
        fmt.Printf("[✓] Loaded graph with %d nodes and %d edges\n", len(graphData.Nodes), len(graphData.Links))
    } else {
        // Create sample data
        graphData = createSampleGraph()
        fmt.Printf("[✓] Created sample graph with %d nodes and %d edges\n", len(graphData.Nodes), len(graphData.Links))
    }
    
    // Start HTTP server
    fmt.Printf("[*] Starting server on port %d\n", port)
    
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        renderGraph(w, graphData)
    })
    
    http.HandleFunc("/api/graph", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(graphData)
    })
    
    if openBrowser {
        url := fmt.Sprintf("http://localhost:%d", port)
        fmt.Printf("[*] Opening browser at %s\n", url)
        openBrowserURL(url)
    }
    
    fmt.Printf("\n[✓] Server running at http://localhost:%d\n", port)
    fmt.Println("Press Ctrl+C to stop")
    
    http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func createSampleGraph() GraphData {
    return GraphData{
        Nodes: []Node{
            {ID: "sub1", Label: "api.example.com", Group: "subdomain", Size: 10},
            {ID: "sub2", Label: "admin.example.com", Group: "subdomain", Size: 8},
            {ID: "sub3", Label: "cdn.example.com", Group: "subdomain", Size: 6},
            {ID: "api1", Label: "/api/users", Group: "api", Size: 5},
            {ID: "api2", Label: "/api/auth", Group: "api", Size: 5},
            {ID: "vuln1", Label: "SQL Injection", Group: "vulnerability", Size: 7},
            {ID: "vuln2", Label: "XSS", Group: "vulnerability", Size: 6},
            {ID: "tech1", Label: "Nginx", Group: "technology", Size: 4},
            {ID: "tech2", Label: "PostgreSQL", Group: "technology", Size: 4},
        },
        Links: []Link{
            {Source: "sub1", Target: "api1", Type: "hosts"},
            {Source: "sub1", Target: "api2", Type: "hosts"},
            {Source: "sub2", Target: "vuln1", Type: "vulnerable_to"},
            {Source: "sub1", Target: "tech1", Type: "uses"},
            {Source: "sub1", Target: "tech2", Type: "uses"},
            {Source: "api1", Target: "vuln2", Type: "vulnerable_to"},
        },
    }
}

func renderGraph(w http.ResponseWriter, data GraphData) {
    tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>ReconForge Graph Viewer</title>
    <script src="https://d3js.org/d3.v7.min.js"></script>
    <style>
        body { margin: 0; font-family: Arial, sans-serif; }
        #graph-container { width: 100vw; height: 100vh; background: #1a1a2e; }
        .node circle { stroke: #fff; stroke-width: 2px; }
        .node text { fill: #fff; font-size: 12px; }
        .link { stroke: #999; stroke-opacity: 0.6; }
        .tooltip { position: absolute; background: rgba(0,0,0,0.8); color: white; padding: 5px 10px; border-radius: 5px; font-size: 12px; pointer-events: none; }
        .controls { position: absolute; top: 20px; right: 20px; background: rgba(0,0,0,0.7); padding: 10px; border-radius: 5px; color: white; z-index: 100; }
        button { margin: 5px; padding: 5px 10px; cursor: pointer; }
        .legend { position: absolute; bottom: 20px; left: 20px; background: rgba(0,0,0,0.7); padding: 10px; border-radius: 5px; color: white; font-size: 12px; }
        .legend-item { display: inline-block; margin-right: 15px; }
        .legend-color { display: inline-block; width: 12px; height: 12px; border-radius: 50%; margin-right: 5px; }
    </style>
</head>
<body>
    <div id="graph-container"></div>
    <div class="controls">
        <button onclick="zoomIn()">Zoom In</button>
        <button onclick="zoomOut()">Zoom Out</button>
        <button onclick="resetZoom()">Reset</button>
        <button onclick="centerGraph()">Center</button>
    </div>
    <div class="legend">
        <div class="legend-item"><span class="legend-color" style="background: #ff6b6b;"></span> Subdomain</div>
        <div class="legend-item"><span class="legend-color" style="background: #4ecdc4;"></span> API</div>
        <div class="legend-item"><span class="legend-color" style="background: #ffe66d;"></span> Vulnerability</div>
        <div class="legend-item"><span class="legend-color" style="background: #95e77e;"></span> Technology</div>
    </div>
    <div class="tooltip" style="display:none"></div>
    <div id="graph-data" style="display:none">{{.}}</div>
    <script>
        const graphData = JSON.parse(document.getElementById('graph-data').innerText);
        
        const width = window.innerWidth;
        const height = window.innerHeight;
        
        const svg = d3.select("#graph-container")
            .append("svg")
            .attr("width", width)
            .attr("height", height)
            .call(d3.zoom().on("zoom", (event) => {
                g.attr("transform", event.transform);
            }))
            .append("g");
        
        const color = d3.scaleOrdinal()
            .domain(["subdomain", "api", "vulnerability", "technology"])
            .range(["#ff6b6b", "#4ecdc4", "#ffe66d", "#95e77e"]);
        
        const simulation = d3.forceSimulation(graphData.nodes)
            .force("link", d3.forceLink(graphData.links).id(d => d.id).distance(150))
            .force("charge", d3.forceManyBody().strength(-300))
            .force("center", d3.forceCenter(width / 2, height / 2));
        
        const link = svg.append("g")
            .selectAll("line")
            .data(graphData.links)
            .enter()
            .append("line")
            .attr("class", "link")
            .style("stroke", "#999")
            .style("stroke-width", 2);
        
        const node = svg.append("g")
            .selectAll("g")
            .data(graphData.nodes)
            .enter()
            .append("g")
            .call(d3.drag()
                .on("start", dragstarted)
                .on("drag", dragged)
                .on("end", dragended));
        
        node.append("circle")
            .attr("r", d => d.size || 8)
            .attr("fill", d => color(d.group))
            .attr("stroke", "#fff")
            .attr("stroke-width", 2);
        
        node.append("text")
            .attr("dx", 12)
            .attr("dy", 4)
            .text(d => d.label)
            .style("fill", "#fff")
            .style("font-size", "12px");
        
        node.on("mouseover", function(event, d) {
            d3.select(".tooltip")
                .style("display", "block")
                .html(`<strong>${d.label}</strong><br>Type: ${d.group}<br>ID: ${d.id}`)
                .style("left", (event.pageX + 10) + "px")
                .style("top", (event.pageY - 10) + "px");
        }).on("mouseout", function() {
            d3.select(".tooltip").style("display", "none");
        });
        
        simulation.on("tick", () => {
            link
                .attr("x1", d => d.source.x)
                .attr("y1", d => d.source.y)
                .attr("x2", d => d.target.x)
                .attr("y2", d => d.target.y);
            
            node
                .attr("transform", d => `translate(${d.x}, ${d.y})`);
        });
        
        function dragstarted(event, d) {
            if (!event.active) simulation.alphaTarget(0.3).restart();
            d.fx = d.x;
            d.fy = d.y;
        }
        
        function dragged(event, d) {
            d.fx = event.x;
            d.fy = event.y;
        }
        
        function dragended(event, d) {
            if (!event.active) simulation.alphaTarget(0);
            d.fx = null;
            d.fy = null;
        }
        
        let zoom = d3.zoom().scaleExtent([0.1, 4]);
        let g = svg;
        
        function zoomIn() { svg.transition().call(zoom.scaleBy, 1.2); }
        function zoomOut() { svg.transition().call(zoom.scaleBy, 0.8); }
        function resetZoom() { svg.transition().call(zoom.transform, d3.zoomIdentity); }
        function centerGraph() { svg.transition().call(zoom.transform, d3.zoomIdentity); }
        
        window.zoomIn = zoomIn;
        window.zoomOut = zoomOut;
        window.resetZoom = resetZoom;
        window.centerGraph = centerGraph;
    </script>
</body>
</html>`
    
    t := template.Must(template.New("graph").Parse(tmpl))
    
    // Convert graph data to JSON string
    dataJSON, _ := json.Marshal(graphData)
    t.Execute(w, template.JS(dataJSON))
}

func openBrowserURL(url string) {
    var err error
    switch os.Getenv("OSTYPE") {
    case "linux-gnu":
        err = exec.Command("xdg-open", url).Start()
    case "darwin":
        err = exec.Command("open", url).Start()
    default:
        err = exec.Command("xdg-open", url).Start()
    }
    if err != nil {
        fmt.Printf("[!] Could not open browser: %v\n", err)
    }
}
