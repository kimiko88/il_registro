<template>
  <Line :data="chartData" :options="chartOptions" />
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
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

const { t } = useI18n();

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
  data: {
    type: Array,
    default: () => []
  }
});

const chartData = computed(() => ({
  labels: (props.data || []).map(d => d.timestamp),
  datasets: [{
    label: t('admin.apiRequests') || 'API Requests',
    backgroundColor: '#f87979',
    data: (props.data || []).map(d => d.value)
  }]
}));

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false
};
</script>
