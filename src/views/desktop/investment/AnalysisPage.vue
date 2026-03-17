<template>
    <div class="page-content">
        <div class="page-header">
            <h1 class="page-title">{{ tt('Performance Analysis') }}</h1>
        </div>
        <div class="page-body">
            <v-card class="mb-4">
                <v-card-title>{{ tt('Performance Metrics') }}</v-card-title>
                <v-card-text>
                    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <v-card variant="outlined" class="p-4">
                            <div class="text-center">
                                <div class="text-2xl font-bold">5.0%</div>
                                <div class="text-gray-500">{{ tt('Annual Return') }}</div>
                            </div>
                        </v-card>
                        <v-card variant="outlined" class="p-4">
                            <div class="text-center">
                                <div class="text-2xl font-bold">12.0%</div>
                                <div class="text-gray-500">{{ tt('Volatility') }}</div>
                            </div>
                        </v-card>
                        <v-card variant="outlined" class="p-4">
                            <div class="text-center">
                                <div class="text-2xl font-bold">0.85</div>
                                <div class="text-gray-500">{{ tt('Sharpe Ratio') }}</div>
                            </div>
                        </v-card>
                    </div>
                </v-card-text>
            </v-card>
            <v-card class="mb-4">
                <v-card-title>{{ tt('Performance Over Time') }}</v-card-title>
                <v-card-text>
                    <div class="h-80">
                        <!-- 时间表现图表 -->
                        <div class="flex items-center justify-center h-full">
                            <div class="text-gray-500">{{ tt('Performance over time chart will be displayed here') }}</div>
                        </div>
                    </div>
                </v-card-text>
            </v-card>
            <v-card class="mb-4">
                <v-card-title>{{ tt('Asset Class Performance') }}</v-card-title>
                <v-card-text>
                    <div class="h-80">
                        <!-- 资产类别表现图表 -->
                        <div class="flex items-center justify-center h-full">
                            <div class="text-gray-500">{{ tt('Asset class performance chart will be displayed here') }}</div>
                        </div>
                    </div>
                </v-card-text>
            </v-card>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales/helpers.ts';

const { tt } = useI18n();
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