<template>
    <div class="page-content">
        <div class="page-header">
            <h1 class="page-title">{{ tt('Asset Management') }}</h1>
        </div>
        <div class="page-body">
            <v-card class="mb-4">
                <v-card-title>{{ tt('Asset List') }}</v-card-title>
                <v-card-text>
                    <div class="overflow-x-auto">
                        <table class="w-full">
                            <thead>
                                <tr>
                                    <th class="text-left">Asset</th>
                                    <th class="text-left">Type</th>
                                    <th class="text-left">Symbol</th>
                                    <th class="text-right">Current Price</th>
                                    <th class="text-right">Holdings</th>
                                    <th class="text-right">Total Value</th>
                                    <th class="text-right">Action</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="(asset, index) in assets" :key="index">
                                    <td>{{ asset.name }}</td>
                                    <td>{{ asset.type }}</td>
                                    <td>{{ asset.symbol }}</td>
                                    <td class="text-right">{{ asset.price }}</td>
                                    <td class="text-right">{{ asset.holdings }}</td>
                                    <td class="text-right">{{ asset.value }}</td>
                                    <td class="text-right">
                                        <v-btn color="primary" variant="text" size="small">{{ tt('Edit') }}</v-btn>
                                        <v-btn color="error" variant="text" size="small">{{ tt('Delete') }}</v-btn>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </v-card-text>
                <v-card-actions>
                    <v-btn color="primary">{{ tt('Add Asset') }}</v-btn>
                </v-card-actions>
            </v-card>
            <v-card class="mb-4">
                <v-card-title>{{ tt('Asset Categories') }}</v-card-title>
                <v-card-text>
                    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <v-card variant="outlined" class="p-4">
                            <div class="text-center">
                                <div class="text-2xl font-bold">$58,000</div>
                                <div class="text-gray-500">Stocks</div>
                                <div class="text-sm text-gray-500">58%</div>
                            </div>
                        </v-card>
                        <v-card variant="outlined" class="p-4">
                            <div class="text-center">
                                <div class="text-2xl font-bold">$50,000</div>
                                <div class="text-gray-500">ETFs</div>
                                <div class="text-sm text-gray-500">50%</div>
                            </div>
                        </v-card>
                        <v-card variant="outlined" class="p-4">
                            <div class="text-center">
                                <div class="text-2xl font-bold">$4,500</div>
                                <div class="text-gray-500">Bonds</div>
                                <div class="text-sm text-gray-500">4.5%</div>
                            </div>
                        </v-card>
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
    { name: 'Apple Inc.', type: 'Stock', symbol: 'AAPL', price: '$180.00', holdings: '100', value: '$18,000' },
    { name: 'Microsoft Corporation', type: 'Stock', symbol: 'MSFT', price: '$400.00', holdings: '50', value: '$20,000' },
    { name: 'Alphabet Inc.', type: 'Stock', symbol: 'GOOGL', price: '$130.00', holdings: '30', value: '$3,900' },
    { name: 'SPDR S&P 500 ETF Trust', type: 'ETF', symbol: 'SPY', price: '$500.00', holdings: '100', value: '$50,000' },
    { name: 'iShares 20+ Year Treasury Bond ETF', type: 'Bond', symbol: 'TLT', price: '$90.00', holdings: '50', value: '$4,500' },
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