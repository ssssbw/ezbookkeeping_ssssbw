<template>
    <div class="page-content">
        <div class="page-header">
            <h1 class="page-title">{{ tt('Strategy Configuration') }}</h1>
        </div>
        <div class="page-body">
            <v-card>
                <v-card-title class="text-subtitle-1">{{ tt('Target Allocation') }}</v-card-title>
                <v-card-text>
                    <v-row>
                        <v-col cols="12" md="3" v-for="cat in categoryTargets" :key="cat.key">
                            <v-text-field
                                v-model.number="cat.target"
                                :label="cat.label"
                                type="number"
                                min="0"
                                max="100"
                                suffix="%"
                                variant="outlined"
                                density="compact"
                            />
                        </v-col>
                    </v-row>
                    <div class="d-flex align-center mt-2">
                        <span class="text-body-2 me-2">{{ tt('Total') }}:</span>
                        <span class="text-body-2 font-weight-bold" :class="totalTarget === 100 ? 'text-success' : 'text-error'">
                            {{ totalTarget }}%
                        </span>
                        <span v-if="totalTarget !== 100" class="text-caption text-error ms-2">
                            ({{ totalTarget > 100 ? '+' : '' }}{{ totalTarget - 100 }}%)
                        </span>
                    </div>
                </v-card-text>
            </v-card>

            <v-card>
                <v-card-title class="text-subtitle-1">{{ tt('Strategy Analysis') }}</v-card-title>
                <v-card-text class="pa-0">
                    <v-data-table
                        :headers="headers"
                        :items="analysisData"
                        :hover="true"
                        :items-per-page="-1"
                        class="strategy-table"
                    >
                        <template #item.currentPct="{ item }">
                            <span class="text-body-2">{{ item.currentPct.toFixed(1) }}%</span>
                        </template>
                        <template #item.targetPct="{ item }">
                            <span class="text-body-2">{{ item.targetPct.toFixed(1) }}%</span>
                        </template>
                        <template #item.deviation="{ item }">
                            <span class="text-body-2" :class="getDeviationClass(item.deviation)">
                                {{ item.deviation >= 0 ? '+' : '' }}{{ item.deviation.toFixed(1) }}%
                            </span>
                        </template>
                        <template #item.action="{ item }">
                            <v-chip size="small" :color="getActionColor(item.deviation)" variant="tonal">
                                {{ item.deviation > 0 ? tt('Reduce') : item.deviation < 0 ? tt('Increase') : '--' }}
                            </v-chip>
                        </template>
                        <template #no-data>
                            <div class="text-center py-6 text-medium-emphasis">
                                {{ tt('No holdings data') }}
                            </div>
                        </template>
                    </v-data-table>
                </v-card-text>
            </v-card>

            <v-row>
                <v-col cols="12" md="6">
                    <v-card>
                        <v-card-title class="text-subtitle-1">{{ tt('Current Allocation') }}</v-card-title>
                        <v-card-text>
                            <v-chart v-if="currentAllocationData.length > 0" autoresize class="allocation-chart" :option="currentChartOptions" />
                            <div v-else class="text-center py-6 text-medium-emphasis">
                                {{ tt('No holdings data') }}
                            </div>
                        </v-card-text>
                    </v-card>
                </v-col>
                <v-col cols="12" md="6">
                    <v-card>
                        <v-card-title class="text-subtitle-1">{{ tt('Target Allocation') }}</v-card-title>
                        <v-card-text>
                            <v-chart v-if="targetAllocationData.length > 0" autoresize class="allocation-chart" :option="targetChartOptions" />
                            <div v-else class="text-center py-6 text-medium-emphasis">
                                {{ tt('No holdings data') }}
                            </div>
                        </v-card-text>
                    </v-card>
                </v-col>
            </v-row>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useInvestmentStore } from '@/stores/investment.ts';

const { tt } = useI18n();
const investmentStore = useInvestmentStore();

const aggregatedHoldings = computed(() => investmentStore.aggregatedHoldings);

interface CategoryTarget {
    key: string;
    label: string;
    target: number;
}

const categoryTargets = ref<CategoryTarget[]>([
    { key: 'equity', label: tt('Equity'), target: 60 },
    { key: 'fixed_income', label: tt('Fixed Income'), target: 30 },
    { key: 'commodity', label: tt('Commodity'), target: 10 },
    { key: 'digital', label: tt('Digital'), target: 0 },
]);

const totalTarget = computed(() => categoryTargets.value.reduce((sum, c) => sum + c.target, 0));

const totalMarketValue = computed(() => aggregatedHoldings.value.reduce((sum, h) => sum + h.totalMarketValue, 0));

interface AnalysisRow {
    category: string;
    label: string;
    currentValue: number;
    currentPct: number;
    targetPct: number;
    deviation: number;
}

const analysisData = computed<AnalysisRow[]>(() => {
    if (totalMarketValue.value === 0) return [];

    const categoryValues = new Map<string, number>();
    for (const h of aggregatedHoldings.value) {
        const cat = h.category || 'other';
        categoryValues.set(cat, (categoryValues.get(cat) || 0) + h.totalMarketValue);
    }

    return categoryTargets.value.map(t => {
        const current = categoryValues.get(t.key) || 0;
        const currentPct = (current / totalMarketValue.value) * 100;
        return {
            category: t.key,
            label: t.label,
            currentValue: current,
            currentPct,
            targetPct: t.target,
            deviation: currentPct - t.target,
        };
    });
});

const headers = computed(() => [
    { key: 'label', title: tt('Category'), sortable: false },
    { key: 'currentPct', title: tt('Current Allocation'), sortable: false },
    { key: 'targetPct', title: tt('Target Allocation'), sortable: false },
    { key: 'deviation', title: tt('Deviation'), sortable: false },
    { key: 'action', title: tt('Rebalance'), sortable: false },
]);

function getDeviationClass(deviation: number): string {
    if (Math.abs(deviation) < 1) return '';
    return deviation > 0 ? 'text-error' : 'text-success';
}

function getActionColor(deviation: number): string {
    if (Math.abs(deviation) < 1) return 'default';
    return deviation > 0 ? 'error' : 'success';
}

const currentAllocationData = computed(() => {
    if (totalMarketValue.value === 0) return [];
    const categoryValues = new Map<string, number>();
    for (const h of aggregatedHoldings.value) {
        const cat = h.category || 'other';
        categoryValues.set(cat, (categoryValues.get(cat) || 0) + h.totalMarketValue);
    }

    return Array.from(categoryValues.entries()).map(([cat, value]) => ({
        name: formatCategoryName(cat),
        value,
    }));
});

const targetAllocationData = computed(() => {
    return categoryTargets.value
        .filter(t => t.target > 0)
        .map(t => ({
            name: t.label,
            value: t.target,
        }));
});

function formatCategoryName(category: string): string {
    switch (category) {
        case 'equity': return tt('Equity');
        case 'fixed_income': return tt('Fixed Income');
        case 'commodity': return tt('Commodity');
        case 'digital': return tt('Digital');
        default: return category;
    }
}

const PIE_COLORS = ['#c67e48', '#4caf50', '#ff9800', '#9c27b0', '#607d8b'];

function makePieOption(data: { name: string; value: number }[]) {
    return {
        tooltip: { trigger: 'item', formatter: '{b}: {d}%' },
        series: [
            {
                type: 'pie',
                radius: ['40%', '70%'],
                avoidLabelOverlap: true,
                itemStyle: { borderRadius: 4, borderColor: '#fff', borderWidth: 2 },
                label: { show: false },
                emphasis: { label: { show: true, fontSize: 14, fontWeight: 'bold' } },
                data: data.map((item, i) => ({
                    value: item.value,
                    name: item.name,
                    itemStyle: { color: PIE_COLORS[i % PIE_COLORS.length] },
                })),
            },
        ],
    };
}

const currentChartOptions = computed(() => makePieOption(currentAllocationData.value));
const targetChartOptions = computed(() => makePieOption(targetAllocationData.value));

onMounted(async () => {
    await investmentStore.loadHoldings({ force: false });
});
</script>

<style scoped>
.page-content {
    padding: 24px;
    height: 100%;
    display: flex;
    flex-direction: column;
    overflow-y: auto;
}

.page-header {
    margin-bottom: 24px;
}

.page-title {
    font-size: 24px;
    font-weight: 600;
    margin: 0;
}

.page-body {
    display: flex;
    flex-direction: column;
    gap: 24px;
}

.allocation-chart {
    width: 100%;
    height: 240px;
}

.strategy-table :deep(.v-data-table__td) {
    padding-top: 8px;
    padding-bottom: 8px;
}
</style>
