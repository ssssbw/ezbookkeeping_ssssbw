import { AssetCategory, InvestmentMarket } from '@/models/investment.ts';

type TranslateFn = (key: string) => string;

const DIVISOR = 10000;

export function formatPrice(value: number | undefined | null): string {
    if (value === undefined || value === null) return '--';
    const num = value / DIVISOR;
    if (num >= 1000) return num.toFixed(2);
    if (num >= 1) return num.toFixed(4);
    return num.toFixed(6);
}

export function formatCurrencyValue(value: number | undefined | null, currency: string): string {
    if (value === undefined || value === null) return '--';
    const num = value / DIVISOR;
    const abs = Math.abs(num);
    let formatted: string;
    if (abs >= 1000000) {
        formatted = (abs / 1000000).toFixed(2) + 'M';
    } else if (abs >= 10000) {
        formatted = (abs / 10000).toFixed(2) + 'W';
    } else if (abs >= 1000) {
        formatted = (abs / 1000).toFixed(2) + 'K';
    } else {
        formatted = abs.toFixed(2);
    }
    const sign = num < 0 ? '-' : '';
    return sign + getCurrencySymbol(currency) + formatted;
}

export function formatReturnRate(value: number | undefined | null): string {
    if (value === undefined || value === null) return '--';
    const pct = value / 100;
    const sign = pct >= 0 ? '+' : '';
    return sign + pct.toFixed(2) + '%';
}

export function formatQuantity(value: number | undefined | null): string {
    if (value === undefined || value === null) return '--';
    const num = value / DIVISOR;
    if (num >= 1000000) return (num / 1000000).toFixed(2) + 'M';
    if (num >= 10000) return (num / 10000).toFixed(2) + 'W';
    if (num >= 1000) return (num / 1000).toFixed(2) + 'K';
    return num.toFixed(2);
}

export function formatMarket(market: number, tt: TranslateFn): string {
    switch (market) {
        case InvestmentMarket.CN: return tt('Market CN');
        case InvestmentMarket.HK: return tt('Market HK');
        case InvestmentMarket.US: return tt('Market US');
        default: return market.toString();
    }
}

export function formatCategory(category: string, tt: TranslateFn): string {
    switch (category) {
        case AssetCategory.Equity: return tt('Equity');
        case AssetCategory.FixedIncome: return tt('Fixed Income');
        case AssetCategory.Commodity: return tt('Commodity');
        case AssetCategory.Digital: return tt('Digital');
        default: return category;
    }
}

export function formatIndustry(industry: string, tt: TranslateFn): string {
    const industryMap: Record<string, string> = {
        'technology': tt('Technology'),
        'healthcare': tt('Healthcare'),
        'consumer': tt('Consumer'),
        'finance': tt('Finance'),
        'energy': tt('Energy'),
        'industrial': tt('Industrial'),
        'real_estate': tt('Real Estate'),
        'materials': tt('Materials'),
        'utilities': tt('Utilities'),
        'telecom': tt('Telecom'),
        'other': tt('Other Industry'),
    };
    return industryMap[industry] || industry;
}

export function getCurrencySymbol(currency: string): string {
    switch (currency) {
        case 'CNY': return '¥';
        case 'HKD': return 'HK$';
        case 'USD': return '$';
        default: return '';
    }
}

export function getReturnColorClass(value: number | undefined | null): string {
    if (value === undefined || value === null) return '';
    if (value > 0) return 'text-success';
    if (value < 0) return 'text-error';
    return '';
}

export function getCategoryColor(category: string): string {
    switch (category) {
        case AssetCategory.Equity: return 'primary';
        case AssetCategory.FixedIncome: return 'success';
        case AssetCategory.Commodity: return 'warning';
        case AssetCategory.Digital: return 'purple';
        default: return 'default';
    }
}
