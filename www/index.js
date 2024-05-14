const ct = document.getElementById('myChart');

const resp = await fetch('/data');
if (!resp.ok) console.error("you suck", resp);
const data = await resp.json();

const labels = data.
    map(e => e.created_at).
    map(e => new Date(e).getUTCHours());
const uniqueLabels = [...new Set(labels)];

new Chart(ct, {
    type: 'bar',
    data: {
        labels: uniqueLabels,
        datasets: [{
            label: 'energy levels',
            data: data.map(d => d.value),
            borderWidth: 1
        }]
    },
    options: {
        scales: {
            y: {
                beginAtZero: true
            }
        }
    }
});
