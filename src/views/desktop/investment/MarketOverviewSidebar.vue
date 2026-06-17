<template>
    <div class="market-overview">
        <div class="market-header">
            <span class="text-subtitle-2 font-weight-bold">{{ tt('Global Market Overview') }}</span>
            <router-link to="/investment/overview" class="text-caption text-medium-emphasis market-more">
                {{ tt('More') }} >
            </router-link>
        </div>
        <div class="market-list">
            <div v-for="item in indices" :key="item.name" class="market-row">
                <span class="text-body-2 market-name">{{ item.name }}</span>
                <span class="text-body-2 market-value">{{ item.value }}</span>
                <span class="text-body-2 market-change" :class="item.change >= 0 ? 'text-profit' : 'text-loss'">
                    {{ item.change >= 0 ? '+' : '' }}{{ item.change.toFixed(2) }}%
                </span>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales/helpers.ts';

const { tt } = useI18n();

interface MarketIndex {
    name: string;
    value: string;
    change: number;
}

const indices: MarketIndex[] = [
    { name: '上证指数', value: '3,202.35', change: 0.62 },
    { name: '深证成指', value: '9,735.81', change: 0.48 },
    { name: '创业板指', value: '1,875.98', change: 0.76 },
    { name: '恒生指数', value: '18,942.57', change: -0.15 },
    { name: '标普500', value: '5,312.03', change: 0.30 },
    { name: '纳斯达克', value: '16,823.17', change: 0.34 },
    { name: 'COMEX黄金', value: '2,345.40', change: 0.21 },
    { name: '美元指数', value: '104.35', change: -0.12 },
];
</script>

<style scoped>
.market-overview {
    margin: 0 12px 12px;
    padding: 12px;
    border-radius: 8px;
    background: rgba(var(--v-theme-on-surface), 0.05);
}

.market-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 10px;
}

.market-more {
    text-decoration: none;
    cursor: pointer;
}

.market-more:hover {
    text-decoration: underline;
}

.market-list {
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.market-row {
    display: flex;
    align-items: center;
    padding: 3px 0;
}

.market-name {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: 0.8125rem;
}

.market-value {
    width: 72px;
    text-align: right;
    font-variant-numeric: tabular-nums;
    font-size: 0.8125rem;
}

.market-change {
    width: 56px;
    text-align: right;
    font-variant-numeric: tabular-nums;
    font-size: 0.8125rem;
}
</style>
