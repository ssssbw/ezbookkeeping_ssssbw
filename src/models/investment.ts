export enum AssetCategory {
    Equity = 'equity',
    FixedIncome = 'fixed_income',
    Commodity = 'commodity',
    Digital = 'digital'
}

export enum InvestmentMarket {
    CN = 1,
    HK = 2,
    US = 3
}

export enum InvestmentTransactionType {
    Buy = 1,
    Sell = 2,
    DividendCash = 3,
    DividendReinvest = 4,
    Split = 5,
    ConversionOut = 6,
    ConversionIn = 7
}

export interface AssetInfoResponse {
    readonly id: string;
    readonly code: string;
    readonly market: number;
    readonly name: string;
    readonly category: string;
    readonly currency: string;
    readonly industry?: string;
    readonly tags?: string;
    readonly extraInfo?: string;
}

export interface AssetSearchRequest {
    readonly keyword: string;
    readonly limit?: number;
}

export interface AssetGetRequest {
    readonly id: string;
}

export interface AssetCreateRequest {
    readonly code: string;
    readonly market: number;
    readonly name: string;
    readonly category: string;
    readonly currency: string;
    readonly industry?: string;
    readonly tags?: string;
    readonly extraInfo?: string;
}

export interface UserAssetInfoResponse {
    readonly id: string;
    readonly assetId: string;
    readonly isActive: boolean;
    readonly comment?: string;
    readonly asset?: AssetInfoResponse;
}

export interface UserAssetListRequest {
    readonly is_active?: boolean;
}

export interface UserAssetAddRequest {
    readonly assetId: string;
}

export interface UserAssetRemoveRequest {
    readonly assetId: string;
}

export interface InvestmentTransactionInfoResponse {
    readonly id: string;
    readonly assetId: string;
    readonly accountId: string;
    readonly type: number;
    readonly tradeTime: number;
    readonly confirmTime?: number;
    readonly quantity: number;
    readonly price: number;
    readonly amount: number;
    readonly fee: number;
    readonly relatedTransactionId?: string;
    readonly utcOffset: number;
    readonly comment: string;
    readonly assetName?: string;
    readonly assetCode?: string;
    readonly accountName?: string;
}

export interface InvestmentTransactionListRequest {
    readonly asset_id?: string;
    readonly account_id?: string;
    readonly type?: number;
    readonly start_time?: number;
    readonly end_time?: number;
}

export interface InvestmentTransactionGetRequest {
    readonly id: string;
}

export interface InvestmentTransactionCreateRequest {
    readonly assetId: string;
    readonly accountId: string;
    readonly type: number;
    readonly tradeTime: number;
    readonly confirmTime?: number;
    readonly quantity?: number;
    readonly price?: number;
    readonly amount: number;
    readonly fee?: number;
    readonly relatedTransactionId?: string;
    readonly utcOffset: number;
    readonly comment?: string;
    readonly clientSessionId?: string;
}

export interface InvestmentTransactionModifyRequest {
    readonly id: string;
    readonly assetId: string;
    readonly accountId: string;
    readonly type: number;
    readonly tradeTime: number;
    readonly confirmTime?: number;
    readonly quantity?: number;
    readonly price?: number;
    readonly amount: number;
    readonly fee?: number;
    readonly relatedTransactionId?: string;
    readonly utcOffset: number;
    readonly comment?: string;
}

export interface InvestmentTransactionDeleteRequest {
    readonly id: string;
}

export interface MarketDataInfoResponse {
    readonly assetId: string;
    readonly date: number;
    readonly price: number;
    readonly volume?: number;
}

export interface MarketDataListRequest {
    readonly asset_id: string;
    readonly start_time?: number;
    readonly end_time?: number;
}

export interface MarketDataGetRequest {
    readonly asset_id: string;
    readonly date: number;
}

export interface MarketDataCreateRequest {
    readonly assetId: string;
    readonly date: number;
    readonly price: number;
    readonly volume?: number;
}

export interface MarketDataModifyRequest {
    readonly assetId: string;
    readonly date: number;
    readonly price: number;
    readonly volume?: number;
}

export interface MarketDataInitRequest {
    readonly assetCode: string;
    readonly tradeTime: number;
}

export interface MarketDataInitResponse {
    readonly count: number;
    readonly startTime: number;
    readonly endTime: number;
}

export interface MarketDataEstimateRequest {
    readonly assetCode: string;
}

export class InvestmentAsset {
    public id: string;
    public code: string;
    public market: number;
    public name: string;
    public category: string;
    public currency: string;
    public industry?: string;
    public tags?: string;
    public extraInfo?: string;

    constructor(id: string, code: string, market: number, name: string, category: string, currency: string) {
        this.id = id;
        this.code = code;
        this.market = market;
        this.name = name;
        this.category = category;
        this.currency = currency;
    }

    public static of(response: AssetInfoResponse): InvestmentAsset {
        const asset = new InvestmentAsset(response.id, response.code, response.market, response.name, response.category, response.currency);
        asset.industry = response.industry;
        asset.tags = response.tags;
        asset.extraInfo = response.extraInfo;
        return asset;
    }

    public static ofMulti(responses: AssetInfoResponse[]): InvestmentAsset[] {
        return responses.map(response => InvestmentAsset.of(response));
    }
}

export class InvestmentUserAsset {
    public id: string;
    public assetId: string;
    public isActive: boolean;
    public comment: string;
    public asset?: InvestmentAsset;

    constructor(id: string, assetId: string, isActive: boolean, comment: string) {
        this.id = id;
        this.assetId = assetId;
        this.isActive = isActive;
        this.comment = comment;
    }

    public static of(response: UserAssetInfoResponse): InvestmentUserAsset {
        const userAsset = new InvestmentUserAsset(response.id, response.assetId, response.isActive, response.comment || '');
        if (response.asset) {
            userAsset.asset = InvestmentAsset.of(response.asset);
        }
        return userAsset;
    }

    public static ofMulti(responses: UserAssetInfoResponse[]): InvestmentUserAsset[] {
        return responses.map(response => InvestmentUserAsset.of(response));
    }
}

export class InvestmentTransactionItem {
    public id: string;
    public assetId: string;
    public accountId: string;
    public type: number;
    public tradeTime: number;
    public confirmTime?: number;
    public quantity: number;
    public price: number;
    public amount: number;
    public fee: number;
    public relatedTransactionId?: string;
    public utcOffset: number;
    public comment: string;
    public assetName: string;
    public assetCode: string;
    public accountName: string;

    constructor(id: string, assetId: string, accountId: string, type: number, tradeTime: number, amount: number) {
        this.id = id;
        this.assetId = assetId;
        this.accountId = accountId;
        this.type = type;
        this.tradeTime = tradeTime;
        this.amount = amount;
        this.quantity = 0;
        this.price = 0;
        this.fee = 0;
        this.utcOffset = 0;
        this.comment = '';
        this.assetName = '';
        this.assetCode = '';
        this.accountName = '';
    }

    public static of(response: InvestmentTransactionInfoResponse): InvestmentTransactionItem {
        const item = new InvestmentTransactionItem(response.id, response.assetId, response.accountId, response.type, response.tradeTime, response.amount);
        item.confirmTime = response.confirmTime;
        item.quantity = response.quantity;
        item.price = response.price;
        item.fee = response.fee;
        item.relatedTransactionId = response.relatedTransactionId;
        item.utcOffset = response.utcOffset;
        item.comment = response.comment;
        item.assetName = response.assetName || '';
        item.assetCode = response.assetCode || '';
        item.accountName = response.accountName || '';
        return item;
    }

    public static ofMulti(responses: InvestmentTransactionInfoResponse[]): InvestmentTransactionItem[] {
        return responses.map(response => InvestmentTransactionItem.of(response));
    }

    public toCreateRequest(): InvestmentTransactionCreateRequest {
        return {
            assetId: this.assetId,
            accountId: this.accountId,
            type: this.type,
            tradeTime: this.tradeTime,
            confirmTime: this.confirmTime,
            quantity: this.quantity,
            price: this.price,
            amount: this.amount,
            fee: this.fee,
            relatedTransactionId: this.relatedTransactionId,
            utcOffset: this.utcOffset,
            comment: this.comment
        };
    }

    public toModifyRequest(): InvestmentTransactionModifyRequest {
        return {
            id: this.id,
            assetId: this.assetId,
            accountId: this.accountId,
            type: this.type,
            tradeTime: this.tradeTime,
            confirmTime: this.confirmTime,
            quantity: this.quantity,
            price: this.price,
            amount: this.amount,
            fee: this.fee,
            relatedTransactionId: this.relatedTransactionId,
            utcOffset: this.utcOffset,
            comment: this.comment
        };
    }
}

export class InvestmentMarketDataItem {
    public assetId: string;
    public date: number;
    public price: number;
    public volume?: number;

    constructor(assetId: string, date: number, price: number) {
        this.assetId = assetId;
        this.date = date;
        this.price = price;
    }

    public static of(response: MarketDataInfoResponse): InvestmentMarketDataItem {
        const item = new InvestmentMarketDataItem(response.assetId, response.date, response.price);
        item.volume = response.volume;
        return item;
    }

    public static ofMulti(responses: MarketDataInfoResponse[]): InvestmentMarketDataItem[] {
        return responses.map(response => InvestmentMarketDataItem.of(response));
    }
}

// === Analysis Types ===

export interface InvestmentHoldingResponse {
    readonly assetId: string;
    readonly assetCode: string;
    readonly assetName: string;
    readonly category: string;
    readonly currency: string;
    readonly market: number;
    readonly accountId: string;
    readonly quantity: number;
    readonly avgCostPrice: number;
    readonly totalCost: number;
    readonly currentPrice: number;
    readonly marketValue: number;
    readonly unrealizedPnl: number;
    readonly returnRate: number;
}

export interface InvestmentOverviewResponse {
    readonly totalInvestment: number;
    readonly totalMarketValue: number;
    readonly totalUnrealizedPnl: number;
    readonly totalReturnRate: number;
    readonly allocations: InvestmentAllocationItem[];
}

export interface InvestmentAllocationItem {
    readonly category: string;
    readonly value: number;
    readonly percentage: number;
}

export class InvestmentHolding {
    public assetId: string;
    public assetCode: string;
    public assetName: string;
    public category: string;
    public currency: string;
    public market: number;
    public accountId: string;
    public quantity: number;
    public avgCostPrice: number;
    public totalCost: number;
    public currentPrice: number;
    public marketValue: number;
    public unrealizedPnl: number;
    public returnRate: number;

    constructor(assetId: string, assetCode: string, assetName: string, category: string, currency: string, market: number) {
        this.assetId = assetId;
        this.assetCode = assetCode;
        this.assetName = assetName;
        this.category = category;
        this.currency = currency;
        this.market = market;
        this.accountId = '';
        this.quantity = 0;
        this.avgCostPrice = 0;
        this.totalCost = 0;
        this.currentPrice = 0;
        this.marketValue = 0;
        this.unrealizedPnl = 0;
        this.returnRate = 0;
    }

    public static of(response: InvestmentHoldingResponse): InvestmentHolding {
        const holding = new InvestmentHolding(response.assetId, response.assetCode, response.assetName, response.category, response.currency, response.market);
        holding.accountId = response.accountId;
        holding.quantity = response.quantity;
        holding.avgCostPrice = response.avgCostPrice;
        holding.totalCost = response.totalCost;
        holding.currentPrice = response.currentPrice;
        holding.marketValue = response.marketValue;
        holding.unrealizedPnl = response.unrealizedPnl;
        holding.returnRate = response.returnRate;
        return holding;
    }

    public static ofMulti(responses: InvestmentHoldingResponse[]): InvestmentHolding[] {
        return responses.map(response => InvestmentHolding.of(response));
    }
}

export class InvestmentOverview {
    public totalInvestment: number;
    public totalMarketValue: number;
    public totalUnrealizedPnl: number;
    public totalReturnRate: number;
    public allocations: InvestmentAllocationItem[];

    constructor() {
        this.totalInvestment = 0;
        this.totalMarketValue = 0;
        this.totalUnrealizedPnl = 0;
        this.totalReturnRate = 0;
        this.allocations = [];
    }

    public static of(response: InvestmentOverviewResponse): InvestmentOverview {
        const overview = new InvestmentOverview();
        overview.totalInvestment = response.totalInvestment;
        overview.totalMarketValue = response.totalMarketValue;
        overview.totalUnrealizedPnl = response.totalUnrealizedPnl;
        overview.totalReturnRate = response.totalReturnRate;
        overview.allocations = response.allocations || [];
        return overview;
    }
}
