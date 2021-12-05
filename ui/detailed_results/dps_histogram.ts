import { IndividualSimRequest, IndividualSimResult } from '/tbc/core/proto/api.js';

import { ResultComponent, ResultComponentConfig } from './result_component.js';

declare var Chart: any;

export class DpsHistogram extends ResultComponent {
  constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'dps-histogram-root';
    super(config);
	}

	onSimResult(request: IndividualSimRequest, result: IndividualSimResult) {
		const chartBounds = this.rootElem.getBoundingClientRect();

		this.rootElem.textContent = '';
		const chartCanvas = document.createElement("canvas");
		chartCanvas.height = chartBounds.height;
		chartCanvas.width = chartBounds.width;


		const min = result.playerMetrics!.dpsAvg - result.playerMetrics!.dpsStdev;
		const max = result.playerMetrics!.dpsAvg + result.playerMetrics!.dpsStdev;
		const vals: Array<number> = [];
		const colors: Array<string> = [];

		const labels = Object.keys(result.playerMetrics!.dpsHist);
		labels.forEach((k, i) => {
			vals.push(result.playerMetrics!.dpsHist[Number(k)]);
			const val = parseInt(k);
			if (val > min && val < max) {
				colors.push('#1E87F0');
			} else {
				colors.push('#FF6961');
			}
		});

		const ctx = chartCanvas.getContext('2d');
		const chart = new Chart(ctx, {
			type: 'bar',
			data: {
				labels: labels,
				datasets: [{
					data: vals,
					backgroundColor: colors,
				}],
			},
			options: {
				plugins: {
					title: {
						display: true,
						text: 'DPS Histogram',
					},
					legend: {
						display: false,
						labels: {},
					}
				},
				scales: {
					y: {
						beginAtZero: true,
						ticks: {
							display: false
						},
					},
				},
			},
		});
		this.rootElem.appendChild(chartCanvas);
	}
}