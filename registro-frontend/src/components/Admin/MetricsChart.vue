<template>
  <Line :data="chartData" :options="chartOptions" />
</template>

<script setup>
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
} from 'chart.js';
import { Line } from 'vue-chartjs';
import { computed } from 'vue';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

const props = defineProps({
  data: Array
});

const chartData = computed(() => ({
  labels: props.data.map(d => d.timestamp),
  datasets: [{
    label: 'API Requests',
    backgroundColor: '#f87979',
    data: props.data.map(d => d.value)
  }]
}));

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false
};
</script>
