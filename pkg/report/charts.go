package report

import (
    "fmt"
    "strings"
)

type ChartData struct {
    Labels   []string `json:"labels"`
    Datasets []Dataset `json:"datasets"`
}

type Dataset struct {
    Label string   `json:"label"`
    Data  []int    `json:"data"`
    Color string   `json:"backgroundColor"`
}

func GenerateSeverityChart(critical, high, medium, low int) string {
    return fmt.Sprintf(`<canvas id="severityChart" width="400" height="200"></canvas>
<script>
const ctx = document.getElementById('severityChart').getContext('2d');
new Chart(ctx, {
    type: 'bar',
    data: {
        labels: ['Critical', 'High', 'Medium', 'Low'],
        datasets: [{
            label: 'Vulnerabilities by Severity',
            data: [%d, %d, %d, %d],
            backgroundColor: ['#ff4444', '#ff8844', '#ffcc44', '#44ff88']
        }]
    }
});
</script>`, critical, high, medium, low)
}

func GeneratePieChart(data map[string]int) string {
    var labels, values []string
    for label, value := range data {
        labels = append(labels, fmt.Sprintf("'%s'", label))
        values = append(values, fmt.Sprintf("%d", value))
    }
    
    return fmt.Sprintf(`<canvas id="pieChart" width="400" height="200"></canvas>
<script>
const pieCtx = document.getElementById('pieChart').getContext('2d');
new Chart(pieCtx, {
    type: 'pie',
    data: {
        labels: [%s],
        datasets: [{
            data: [%s],
            backgroundColor: ['#00ff88', '#ff4444', '#ff8844', '#ffcc44', '#44ff88']
        }]
    }
});
</script>`, strings.Join(labels, ","), strings.Join(values, ","))
}

func GenerateLineChart(data map[string][]int, labels []string) string {
    // Simplified line chart generation
    return `<canvas id="lineChart" width="400" height="200"></canvas>
<script>
const lineCtx = document.getElementById('lineChart').getContext('2d');
new Chart(lineCtx, {
    type: 'line',
    data: {
        labels: ['Jan', 'Feb', 'Mar', 'Apr'],
        datasets: [{
            label: 'Vulnerabilities',
            data: [5, 8, 3, 10],
            borderColor: '#00ff88'
        }]
    }
});
</script>`
}
