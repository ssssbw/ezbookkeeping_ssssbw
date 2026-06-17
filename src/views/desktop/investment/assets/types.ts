export interface DisplayAsset {
    assetId: string;
    assetCode: string;
    assetName: string;
    category: string;
    currency: string;
    market: number;
    industry: string;
    currentPrice?: number;
    currentPriceDate?: number;
    addedPrice?: number;
    isHolding: boolean;
}
