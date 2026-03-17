<template>
    <div class="page-content">
        <div class="page-header">
            <h1 class="page-title">{{ tt('Investment Portfolio') }}</h1>
        </div>
        <div class="page-body">
            <v-card class="mb-4">
                <v-card-title>{{ tt('Investment Assets') }}</v-card-title>
                <v-card-text>
                    <div class="overflow-x-auto">
                        <table class="w-full">
                            <thead>
                                <tr>
                                    <th class="text-left">Asset</th>
                                    <th class="text-left">Type</th>
                                    <th class="text-right">Quantity</th>
                                    <th class="text-right">Price</th>
                                    <th class="text-right">Value</th>
                                    <th class="text-right">Return</th>
                                    <th class="text-right">Return Rate</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="(asset, index) in assets" :key="index">
                                    <td>{{ asset.name }}</td>
                                    <td>{{ asset.type }}</td>
                                    <td class="text-right">{{ asset.quantity }}</td>
                                    <td class="text-right">{{ asset.price }}</td>
                                    <td class="text-right">{{ asset.value }}</td>
                                    <td class="text-right" :class="asset.return.startsWith('+') ? 'text-green-600' : 'text-red-600'">{{ asset.return }}</td>
                                    <td class="text-right" :class="asset.returnRate.startsWith('+') ? 'text-green-600' : 'text-red-600'">{{ asset.returnRate }}</td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </v-card-text>
            </v-card>
            <v-card class="mb-4">
                <v-card-title>{{ tt('Portfolio Performance') }}</v-card-title>
                <v-card-text>
                    <div class="h-80">
                        <!-- 投资组合表现图表 -->
                        <div class="flex items-center justify-center h-full">
                            <div class="text-gray-500">{{ tt('Portfolio performance chart will be displayed here') }}</div>
                        </div>
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

const assets = ref([
    { name: 'AAPL', type: 'Stock', quantity: '100', price: '$180.00', value: '$18,000', return: '+$1,000', returnRate: '+5.88%' },
    { name: 'MSFT', type: 'Stock', quantity: '50', price: '$400.00', value: '$20,000', return: '+$2,000', returnRate: '+11.11%' },
    { name: 'GOOGL', type: 'Stock', quantity: '30', price: '$130.00', value: '$3,900', return: '-$100', returnRate: '-2.50%' },
    { name: 'SPY', type: 'ETF', quantity: '100', price: '$500.00', value: '$50,000', return: '+$2,500', returnRate: '+5.26%' },
    { name: 'TLT', type: 'Bond', quantity: '50', price: '$90.00', value: '$4,500', return: '+$200', returnRate: '+4.65%' },
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