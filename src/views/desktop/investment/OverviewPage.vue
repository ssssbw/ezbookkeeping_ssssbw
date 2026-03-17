<template>
    <div class="page-content">
        <div class="page-header">
            <h1 class="page-title">{{ tt('Investment Overview') }}</h1>
        </div>
        <div class="page-body">
            <v-card class="mb-4">
                <v-card-title>{{ tt('Investment Portfolio Summary') }}</v-card-title>
                <v-card-text>
                    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <v-card variant="outlined" class="p-4">
                            <div class="text-center">
                                <div class="text-2xl font-bold">$100,000</div>
                                <div class="text-gray-500">{{ tt('Total Investment') }}</div>
                            </div>
                        </v-card>
                        <v-card variant="outlined" class="p-4">
                            <div class="text-center">
                                <div class="text-2xl font-bold text-green-600">+$5,000</div>
                                <div class="text-gray-500">{{ tt('Total Return') }}</div>
                            </div>
                        </v-card>
                        <v-card variant="outlined" class="p-4">
                            <div class="text-center">
                                <div class="text-2xl font-bold text-green-600">5.0%</div>
                                <div class="text-gray-500">{{ tt('Return Rate') }}</div>
                            </div>
                        </v-card>
                    </div>
                </v-card-text>
            </v-card>
            <v-card class="mb-4">
                <v-card-title>{{ tt('Asset Allocation') }}</v-card-title>
                <v-card-text>
                    <div class="h-80">
                        <!-- 资产配置图表 -->
                        <div class="flex items-center justify-center h-full">
                            <div class="text-gray-500">{{ tt('Asset allocation chart will be displayed here') }}</div>
                        </div>
                    </div>
                </v-card-text>
            </v-card>
            <v-card class="mb-4">
                <v-card-title>{{ tt('Recent Transactions') }}</v-card-title>
                <v-card-text>
                    <div class="overflow-x-auto">
                        <table class="w-full">
                            <thead>
                                <tr>
                                    <th class="text-left">Date</th>
                                    <th class="text-left">Type</th>
                                    <th class="text-left">Asset</th>
                                    <th class="text-right">Amount</th>
                                    <th class="text-right">Status</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="(transaction, index) in recentTransactions" :key="index">
                                    <td>{{ transaction.date }}</td>
                                    <td>{{ transaction.type }}</td>
                                    <td>{{ transaction.asset }}</td>
                                    <td class="text-right">{{ transaction.amount }}</td>
                                    <td class="text-right">
                                        <span :class="transaction.status === 'Completed' ? 'text-green-600' : 'text-yellow-600'">{{ transaction.status }}</span>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </v-card-text>
            </v-card>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from '@/locales/helpers.ts';

const { tt } = useI18n();

const recentTransactions = ref([
    { date: '2026-03-15', type: 'Buy', asset: 'AAPL', amount: '$1,000', status: 'Completed' },
    { date: '2026-03-14', type: 'Sell', asset: 'MSFT', amount: '$2,500', status: 'Completed' },
    { date: '2026-03-13', type: 'Buy', asset: 'GOOGL', amount: '$1,500', status: 'Pending' },
    { date: '2026-03-12', type: 'Dividend', asset: 'AAPL', amount: '$50', status: 'Completed' },
]);
</script>

<style scoped>
.page-content {
    padding: 24px;
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
</style>