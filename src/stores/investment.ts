import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import {
    InvestmentUserAsset,
    InvestmentTransactionItem,
    InvestmentMarketDataItem,
    InvestmentHolding,
    InvestmentOverview,
    type MarketDataInitResponse
} from '@/models/investment.ts';

import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

export const useInvestmentStore = defineStore('investment', () => {
    const userAssets = ref<InvestmentUserAsset[]>([]);
    const userAssetsMap = ref<Record<string, InvestmentUserAsset>>({});
    const userAssetsStateInvalid = ref<boolean>(true);

    const transactions = ref<InvestmentTransactionItem[]>([]);
    const transactionsStateInvalid = ref<boolean>(true);

    const latestMarketDataMap = ref<Record<string, InvestmentMarketDataItem>>({});

    const holdings = ref<InvestmentHolding[]>([]);
    const holdingsMap = ref<Record<string, InvestmentHolding>>({});
    const holdingsStateInvalid = ref<boolean>(true);

    const overview = ref<InvestmentOverview | null>(null);
    const overviewStateInvalid = ref<boolean>(true);

    const activeUserAssets = computed<InvestmentUserAsset[]>(() => {
        return userAssets.value.filter(ua => ua.isActive);
    });

    const userAssetsCount = computed<number>(() => {
        return userAssets.value.length;
    });

    const holdingsCount = computed<number>(() => {
        return holdings.value.length;
    });

    const holdingsSummary = computed(() => {
        let totalMarketValue = 0;
        let totalCost = 0;
        let totalUnrealizedPnl = 0;

        for (const h of holdings.value) {
            totalMarketValue += h.marketValue;
            totalCost += h.totalCost;
            totalUnrealizedPnl += h.unrealizedPnl;
        }

        const totalReturnRate = totalCost !== 0
            ? Math.round((totalUnrealizedPnl / totalCost) * 10000)
            : 0;

        return { totalMarketValue, totalCost, totalUnrealizedPnl, totalReturnRate };
    });

    interface AggregatedHolding {
        assetId: string;
        assetCode: string;
        assetName: string;
        category: string;
        market: number;
        currency: string;
        totalQuantity: number;
        totalMarketValue: number;
        totalCost: number;
        unrealizedPnl: number;
        weightedReturnRate: number;
        currentPrice: number;
        currentPriceDate: number;
        holdings: InvestmentHolding[];
    }

    const aggregatedHoldings = computed<AggregatedHolding[]>(() => {
        const map = new Map<string, AggregatedHolding>();

        for (const h of holdings.value) {
            let agg = map.get(h.assetId);
            if (!agg) {
                agg = {
                    assetId: h.assetId,
                    assetCode: h.assetCode,
                    assetName: h.assetName,
                    category: h.category,
                    market: h.market,
                    currency: h.currency,
                    totalQuantity: 0,
                    totalMarketValue: 0,
                    totalCost: 0,
                    unrealizedPnl: 0,
                    weightedReturnRate: 0,
                    currentPrice: h.currentPrice,
                    currentPriceDate: h.currentPriceDate,
                    holdings: [],
                };
                map.set(h.assetId, agg);
            }

            agg.totalQuantity += h.quantity;
            agg.totalMarketValue += h.marketValue;
            agg.totalCost += h.totalCost;
            agg.unrealizedPnl += h.unrealizedPnl;
            agg.holdings.push(h);
        }

        for (const agg of map.values()) {
            agg.weightedReturnRate = agg.totalCost !== 0
                ? Math.round((agg.unrealizedPnl / agg.totalCost) * 10000)
                : 0;
        }

        return Array.from(map.values());
    });

    function loadUserAssets({ force }: { force: boolean }): Promise<InvestmentUserAsset[]> {
        if (!force && !userAssetsStateInvalid.value) {
            return new Promise((resolve) => { resolve(userAssets.value); });
        }

        return new Promise((resolve, reject) => {
            services.getUserAssets().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve user asset list' });
                    return;
                }

                if (userAssetsStateInvalid.value) {
                    userAssetsStateInvalid.value = false;
                }

                const list = InvestmentUserAsset.ofMulti(data.result);
                userAssets.value = list;

                const map: Record<string, InvestmentUserAsset> = {};
                for (const ua of list) {
                    map[ua.assetId] = ua;
                }
                userAssetsMap.value = map;

                resolve(list);
            }).catch(error => {
                logger.error('failed to load user asset list', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve user asset list' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function addUserAsset({ assetId }: { assetId: string }): Promise<string> {
        return new Promise((resolve, reject) => {
            services.addUserAsset({ assetId }).then(response => {
                const data = response.data;

                if (!data || !data.success) {
                    reject({ message: 'Unable to add user asset' });
                    return;
                }

                userAssetsStateInvalid.value = true;
                resolve('ok');
            }).catch(error => {
                logger.error('failed to add user asset', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to add user asset' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function removeUserAsset({ assetId }: { assetId: string }): Promise<string> {
        return new Promise((resolve, reject) => {
            services.removeUserAsset({ assetId }).then(response => {
                const data = response.data;

                if (!data || !data.success) {
                    reject({ message: 'Unable to remove user asset' });
                    return;
                }

                userAssetsStateInvalid.value = true;
                resolve('ok');
            }).catch(error => {
                logger.error('failed to remove user asset', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to remove user asset' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function loadTransactions({ assetId, accountId, force }: { assetId?: string; accountId?: string; force: boolean }): Promise<InvestmentTransactionItem[]> {
        if (!force && !transactionsStateInvalid.value && !assetId && !accountId) {
            return new Promise((resolve) => { resolve(transactions.value); });
        }

        return new Promise((resolve, reject) => {
            services.getInvestmentTransactions({
                asset_id: assetId,
                account_id: accountId
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve investment transactions' });
                    return;
                }

                if (transactionsStateInvalid.value) {
                    transactionsStateInvalid.value = false;
                }

                const list = InvestmentTransactionItem.ofMulti(data.result);
                transactions.value = list;
                resolve(list);
            }).catch(error => {
                logger.error('failed to load investment transactions', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve investment transactions' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function addTransaction({ transaction }: { transaction: InvestmentTransactionItem }): Promise<InvestmentTransactionItem> {
        return new Promise((resolve, reject) => {
            services.addInvestmentTransaction(transaction.toCreateRequest()).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to create investment transaction' });
                    return;
                }

                transactionsStateInvalid.value = true;
                resolve(InvestmentTransactionItem.of(data.result));
            }).catch(error => {
                logger.error('failed to create investment transaction', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to create investment transaction' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function modifyTransaction({ transaction }: { transaction: InvestmentTransactionItem }): Promise<InvestmentTransactionItem> {
        return new Promise((resolve, reject) => {
            services.modifyInvestmentTransaction(transaction.toModifyRequest()).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to modify investment transaction' });
                    return;
                }

                transactionsStateInvalid.value = true;
                resolve(InvestmentTransactionItem.of(data.result));
            }).catch(error => {
                logger.error('failed to modify investment transaction', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to modify investment transaction' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function deleteTransaction({ transactionId }: { transactionId: string }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.deleteInvestmentTransaction({ id: transactionId }).then(response => {
                const data = response.data;

                if (!data || !data.success) {
                    reject({ message: 'Unable to delete investment transaction' });
                    return;
                }

                transactionsStateInvalid.value = true;
                resolve(true);
            }).catch(error => {
                logger.error('failed to delete investment transaction', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to delete investment transaction' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function loadLatestMarketData({ assetId }: { assetId: string }): Promise<InvestmentMarketDataItem | null> {
        return new Promise((resolve) => {
            services.getLatestMarketData({ asset_id: assetId, date: 0 }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    resolve(null);
                    return;
                }

                const item = InvestmentMarketDataItem.of(data.result);
                latestMarketDataMap.value[assetId] = item;
                resolve(item);
            }).catch(error => {
                logger.error('failed to load latest market data', error);
                resolve(null);
            });
        });
    }

    function refreshAllMarketData(): Promise<string> {
        return new Promise((resolve, reject) => {
            services.refreshMarketData().then(response => {
                const data = response.data;

                if (!data || !data.success) {
                    reject({ message: 'Unable to refresh market data' });
                    return;
                }

                resolve('ok');
            }).catch(error => {
                logger.error('failed to refresh market data', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to refresh market data' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function initMarketData({ assetCode, tradeTime }: { assetCode: string; tradeTime: number }): Promise<MarketDataInitResponse> {
        return new Promise((resolve, reject) => {
            services.initMarketData({ assetCode, tradeTime }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to init market data' });
                    return;
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to init market data', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to init market data' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function loadHoldings({ force }: { force: boolean }): Promise<InvestmentHolding[]> {
        if (!force && !holdingsStateInvalid.value) {
            return new Promise((resolve) => { resolve(holdings.value); });
        }

        return new Promise((resolve, reject) => {
            services.getInvestmentHoldings().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve investment holdings' });
                    return;
                }

                if (holdingsStateInvalid.value) {
                    holdingsStateInvalid.value = false;
                }

                const list = InvestmentHolding.ofMulti(data.result);
                holdings.value = list;

                const map: Record<string, InvestmentHolding> = {};
                for (const h of list) {
                    map[h.assetId] = h;
                }
                holdingsMap.value = map;

                resolve(list);
            }).catch(error => {
                logger.error('failed to load investment holdings', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve investment holdings' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function loadOverview({ force }: { force: boolean }): Promise<InvestmentOverview> {
        if (!force && !overviewStateInvalid.value && overview.value) {
            return new Promise((resolve) => { resolve(overview.value!); });
        }

        return new Promise((resolve, reject) => {
            services.getInvestmentOverview().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve investment overview' });
                    return;
                }

                if (overviewStateInvalid.value) {
                    overviewStateInvalid.value = false;
                }

                const result = InvestmentOverview.of(data.result);
                overview.value = result;

                resolve(result);
            }).catch(error => {
                logger.error('failed to load investment overview', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve investment overview' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function resetInvestment(): void {
        userAssets.value = [];
        userAssetsMap.value = {};
        userAssetsStateInvalid.value = true;
        transactions.value = [];
        transactionsStateInvalid.value = true;
        latestMarketDataMap.value = {};
        holdings.value = [];
        holdingsMap.value = {};
        holdingsStateInvalid.value = true;
        overview.value = null;
        overviewStateInvalid.value = true;
    }

    return {
        userAssets,
        userAssetsMap,
        userAssetsStateInvalid,
        transactions,
        transactionsStateInvalid,
        latestMarketDataMap,
        activeUserAssets,
        userAssetsCount,
        holdings,
        holdingsMap,
        holdingsStateInvalid,
        holdingsCount,
        holdingsSummary,
        aggregatedHoldings,
        overview,
        overviewStateInvalid,
        loadUserAssets,
        addUserAsset,
        removeUserAsset,
        loadTransactions,
        addTransaction,
        modifyTransaction,
        deleteTransaction,
        loadLatestMarketData,
        refreshAllMarketData,
        initMarketData,
        loadHoldings,
        loadOverview,
        resetInvestment
    };
});
