#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
投资数据转换脚本 — 旧 CSV 格式 → 新投资模块 API

将 ezBookkeeping 旧格式（Transfer/Income 记录基金操作）转换为新投资系统的
InvestmentAsset + InvestmentTransaction API 请求。

使用方法：
  1. 打开下方 CONFIG 区域，填入你的 JWT Token 和投资池 AccountId
  2. 运行：python convert_investment_data.py --dry-run   （预览）
  3. 运行：python convert_investment_data.py --import    （实际导入）

前置条件：
  - 后端已启动（阶段 1.5 已完成）
  - 已通过现有 Account 体系创建投资池子账户（Category=INVESTMENT）
  - 已登录获取 JWT Token
"""

import csv
import json
import re
import os
import sys
import argparse
from datetime import datetime, timedelta
from collections import defaultdict, OrderedDict

# ============================================================
# CONFIG — 请根据实际情况修改
# ============================================================

CONFIG = {
    # 投资池子账户 ID（通过 Account API 创建后填入）
    "account_id": 3822427607018766337,

    # API 基础地址
    "api_base": "http://localhost:8080/api/v1",

    # JWT Token（登录后从浏览器或 API 获取）
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVG9rZW5JZCI6IjI5NDI5OTc0OTIwMjMyNDgwNTciLCJqdGkiOiIzNzcyNTE3NDk3MDM5NzQ5MTIwIiwidXNlcm5hbWUiOiJ0ZXN0IiwidHlwZSI6MSwiaWF0IjoxNzc5OTU3OTI1LCJleHAiOjE3ODI1NDk5MjV9.MYJAlbkDyHZEf_LHkI5TUGu0QLP25rSGDa8Hf4kSMKU",

    # 输入 CSV 文件（相对于脚本所在目录）
    "input_csv": "../test_data.csv",

    # 默认时区偏移（东八区 = +480 分钟）
    "utc_offset": 480,
}

# ============================================================
# 工具函数
# ============================================================

def parse_fund_info(text):
    """
    从旧 Account/Account2 字段提取基金名称和代码。
    例: "广发全球精选股票(QDII)A-[270023]" → ("广发全球精选股票(QDII)A", "270023")
        "支付宝零钱"                          → ("支付宝零钱", None)
    """
    text = text.strip()
    m = re.search(r'(.+)-\[(\d{6})\]\s*$', text)
    if m:
        return m.group(1).strip(), m.group(2)
    return text, None


def parse_time(s):
    """解析 "2026-03-27 14:58:38" → datetime"""
    return datetime.strptime(s.strip(), "%Y-%m-%d %H:%M:%S")


def to_unix(dt):
    """datetime → unix 秒"""
    return int(dt.timestamp())


def next_business_day(dt):
    """T+1 确认日（跳过周末）"""
    d = dt + timedelta(days=1)
    while d.weekday() >= 5:
        d += timedelta(days=1)
    return d


def amt_to_int(s):
    """金额字符串 → int64 (×10000)"""
    return int(round(float(s.strip()) * 10000))


def guess_asset_type(name):
    """根据名称推断资产类型"""
    if "ETF联" in name:
        return 3   # ETF
    if "ETF" in name:
        return 3
    if "债券" in name:
        return 4   # Bond
    if "股票" in name and "指数" not in name:
        return 2   # Stock
    # 混合、指数、LOF、FOF、QDII → 都归 Fund
    return 1       # Fund


TYPE_LABELS = {
    1: "buy",
    2: "sell",
    3: "dividend_cash",
    4: "dividend_reinvest",
    5: "split",
    6: "conversion_out",
    7: "conversion_in",
}

TYPE_LABELS_CN = {
    1: "买入", 2: "卖出", 3: "现金分红", 4: "红利再投",
    5: "拆分", 6: "转出", 7: "转入",
}


# ============================================================
# CSV 解析
# ============================================================

def read_csv(path):
    rows = []
    with open(path, "r", encoding="utf-8") as f:
        reader = csv.DictReader(f)
        for r in reader:
            rows.append(r)
    return rows


def extract_assets(rows):
    """从 CSV 提取去重后的资产列表"""
    assets = OrderedDict()   # code → dict
    warnings = []

    for row in rows:
        for field in ("Account", "Account2"):
            val = row[field].strip()
            if not val or val in ("支付宝零钱", "工资-稳健投资"):
                continue
            name, code = parse_fund_info(val)
            if code is None:
                warnings.append(f"  ⚠ 行 {row.get('Time','?')}: 无代码 → '{name}'")
                continue
            if code in assets:
                continue
            assets[code] = {
                "type": guess_asset_type(name),
                "market": 1,       # CN
                "code": code,
                "name": name,
                "currency": "CNY",
                "comment": "",
            }
    return list(assets.values()), warnings


def convert_transactions(rows, account_id, utc_offset):
    """
    将 CSV 行转换为新系统交易列表。

    返回:
        transactions: list[dict]  — 每条可直接 POST 到 API
        conversion_pairs: list[(out_idx, in_idx)]  — 需要配对的 conversion 索引
        skipped: list[(csv_row_num, reason)]
    """
    transactions = []
    conversion_pairs = []
    skipped = []

    for i, row in enumerate(rows):
        csv_type     = row["Type"].strip()
        category     = row["Category"].strip()
        sub_category = row["Sub Category"].strip()
        account      = row["Account"].strip()
        account2     = row["Account2"].strip()
        amount_s     = row["Amount"].strip()
        time_s       = row["Time"].strip()
        csv_line     = i + 2   # 行号（含表头）

        # ---- 非投资类别跳过 ----
        if category != "金融投资":
            skipped.append((csv_line, "非金融投资类别，跳过"))
            continue

        dt = parse_time(time_s)
        trade_ts = to_unix(dt)
        confirm_ts = to_unix(next_business_day(dt))
        amount_int = amt_to_int(amount_s)

        # =========================================
        # 1) 基金买入
        # =========================================
        if sub_category == "基金买入":
            _, code = parse_fund_info(account2)
            if code is None:
                skipped.append((csv_line, f"买入目标无代码: {account2}"))
                continue
            transactions.append({
                "asset_code": code,
                "accountId": account_id,
                "type": 1,
                "tradeTime": trade_ts,
                "confirmTime": confirm_ts,
                "quantity": 0,    # CSV 无确认份额
                "price": 0,       # CSV 无确认净值
                "amount": amount_int,
                "fee": 0,
                "utcOffset": utc_offset,
                "comment": "",
            })

        # =========================================
        # 2) 基金卖出
        # =========================================
        elif sub_category == "基金卖出":
            _, code = parse_fund_info(account)
            if code is None:
                skipped.append((csv_line, f"卖出源无代码: {account}"))
                continue
            transactions.append({
                "asset_code": code,
                "accountId": account_id,
                "type": 2,
                "tradeTime": trade_ts,
                "confirmTime": confirm_ts,
                "quantity": 0,
                "price": 0,
                "amount": amount_int,
                "fee": 0,
                "utcOffset": utc_offset,
                "comment": "",
            })

        # =========================================
        # 3) 基金转换 → 目标为另一基金
        # =========================================
        elif (sub_category == "基金转换"
              and account2 not in ("支付宝零钱", "工资-稳健投资")):
            name_out, code_out = parse_fund_info(account)
            name_in,  code_in  = parse_fund_info(account2)
            if not code_out or not code_in:
                skipped.append((csv_line, f"转换缺代码: {account} → {account2}"))
                continue

            out_idx = len(transactions)
            transactions.append({
                "asset_code": code_out,
                "accountId": account_id,
                "type": 6,   # conversion_out
                "tradeTime": trade_ts,
                "confirmTime": confirm_ts,
                "quantity": 0,
                "price": 0,
                "amount": amount_int,
                "fee": 0,
                "utcOffset": utc_offset,
                "comment": f"转至 {name_in}",
            })
            in_idx = len(transactions)
            transactions.append({
                "asset_code": code_in,
                "accountId": account_id,
                "type": 7,   # conversion_in
                "tradeTime": trade_ts,
                "confirmTime": confirm_ts,
                "quantity": 0,
                "price": 0,
                "amount": amount_int,
                "fee": 0,
                "utcOffset": utc_offset,
                "comment": f"转自 {name_out}",
            })
            conversion_pairs.append((out_idx, in_idx))

        # =========================================
        # 4) 基金转换 → 退款到零钱/工资池
        # =========================================
        elif (sub_category == "基金转换"
              and account2 in ("支付宝零钱", "工资-稳健投资")):
            _, code = parse_fund_info(account)
            if code is None:
                skipped.append((csv_line, f"退款源无代码: {account}"))
                continue
            transactions.append({
                "asset_code": code,
                "accountId": account_id,
                "type": 2,   # 视为卖出赎回
                "tradeTime": trade_ts,
                "confirmTime": confirm_ts,
                "quantity": 0,
                "price": 0,
                "amount": amount_int,
                "fee": 0,
                "utcOffset": utc_offset,
                "comment": "转换退款",
            })

        # =========================================
        # 5) 分红 / 投资收入
        # =========================================
        elif csv_type == "Income" and sub_category in ("股息分红", "投资收入"):
            _, code = parse_fund_info(account)
            if code is None:
                skipped.append((csv_line, f"分红源无代码: {account}"))
                continue
            transactions.append({
                "asset_code": code,
                "accountId": account_id,
                "type": 3,   # dividend_cash
                "tradeTime": trade_ts,
                "confirmTime": confirm_ts,
                "quantity": 0,
                "price": 0,
                "amount": amount_int,
                "fee": 0,
                "utcOffset": utc_offset,
                "comment": "",
            })

        # =========================================
        # 6) 手续费 — 忽略
        # =========================================
        elif csv_type == "Expense":
            skipped.append((csv_line, "手续费记录已忽略"))

        else:
            skipped.append((csv_line, f"未识别: {csv_type}/{sub_category}"))

    return transactions, conversion_pairs, skipped


# ============================================================
# API 调用（--import 模式）
# ============================================================

def post_json(session, url, payload, token):
    """使用 urllib 发送 POST 请求"""
    import urllib.request
    data = json.dumps(payload, ensure_ascii=False).encode("utf-8")
    req = urllib.request.Request(
        url,
        data=data,
        headers={
            "Content-Type": "application/json",
            "Authorization": f"Bearer {token}",
        },
        method="POST",
    )
    try:
        with urllib.request.urlopen(req) as resp:
            body = json.loads(resp.read().decode("utf-8"))
            return body
    except urllib.error.HTTPError as e:
        body = e.read().decode("utf-8", errors="replace")
        print(f"  ❌ HTTP {e.code}: {body[:200]}")
        return None
    except Exception as e:
        print(f"  ❌ 请求失败: {e}")
        return None


# ============================================================
# 主函数
# ============================================================

def main():
    parser = argparse.ArgumentParser(description="旧投资数据 → 新投资模块 转换工具")
    parser.add_argument("--dry-run", action="store_true", default=True,
                        help="仅预览，不实际调用 API（默认）")
    parser.add_argument("--import", dest="do_import", action="store_true",
                        help="实际调用 API 导入数据")
    parser.add_argument("--csv", default=None,
                        help="覆盖输入 CSV 路径")
    args = parser.parse_args()

    script_dir = os.path.dirname(os.path.abspath(__file__))
    csv_path = args.csv or os.path.join(script_dir, CONFIG["input_csv"])
    out_dir = script_dir  # 输出文件放在 scripts/ 目录

    # ---- 读取 CSV ----
    print(f"📂 读取: {csv_path}")
    rows = read_csv(csv_path)
    print(f"   共 {len(rows)} 行数据\n")

    # ---- 提取资产 ----
    assets, asset_warnings = extract_assets(rows)
    print("=" * 60)
    print(f"📊 提取到 {len(assets)} 个独立资产")
    print("=" * 60)
    for a in assets:
        print(f"   [{a['code']}] {a['name']}  (type={a['type']})")
    if asset_warnings:
        print()
        for w in asset_warnings:
            print(w)

    # ---- 转换交易 ----
    transactions, conv_pairs, skipped = convert_transactions(
        rows, CONFIG["account_id"], CONFIG["utc_offset"])

    # 统计
    type_counts = defaultdict(int)
    for tx in transactions:
        type_counts[tx["type"]] += 1
    total_amount = sum(tx["amount"] for tx in transactions if tx["type"] in (1,))

    print(f"\n{'=' * 60}")
    print(f"📋 交易转换统计")
    print(f"{'=' * 60}")
    for t in sorted(type_counts):
        print(f"   {TYPE_LABELS_CN.get(t, '?'):8s} ({TYPE_LABELS.get(t,'?'):20s}): {type_counts[t]:3d} 条")
    print(f"   {'总计':8s} ({'':20s}): {len(transactions):3d} 条")
    print(f"   买入总金额: {total_amount / 10000:,.2f} CNY")
    print(f"   转换配对数: {len(conv_pairs)} 对")

    if skipped:
        print(f"\n⏭ 跳过 {len(skipped)} 条：")
        for line_no, reason in skipped:
            print(f"   行 {line_no:3d}: {reason}")

    # ---- 生成输出文件 ----
    # 1) assets.json
    assets_file = os.path.join(out_dir, "output_assets.json")
    with open(assets_file, "w", encoding="utf-8") as f:
        json.dump(assets, f, ensure_ascii=False, indent=2)
    print(f"\n✅ 资产列表 → {assets_file}")

    # 2) transactions.json
    # 把 asset_code 转为占位符 assetId，附带 code 用于后续关联
    tx_output = []
    for tx in transactions:
        entry = dict(tx)
        entry["_asset_code"] = entry.pop("asset_code")
        entry["assetId"] = "PLACEHOLDER"  # 导入时替换
        tx_output.append(entry)

    tx_file = os.path.join(out_dir, "output_transactions.json")
    with open(tx_file, "w", encoding="utf-8") as f:
        out = {
            "transactions": tx_output,
            "conversion_pairs": conv_pairs,
            "meta": {
                "account_id": CONFIG["account_id"],
                "utc_offset": CONFIG["utc_offset"],
            },
        }
        json.dump(out, f, ensure_ascii=False, indent=2)
    print(f"✅ 交易列表 → {tx_file}")

    # 3) 生成导入指引
    guide_file = os.path.join(out_dir, "import_guide.md")
    generate_import_guide(guide_file, assets, transactions, conv_pairs, skipped)
    print(f"✅ 导入指引 → {guide_file}")

    # ---- 实际导入模式 ----
    if args.do_import:
        do_import(assets, transactions, conv_pairs)

    # ---- 最终提示 ----
    print(f"""
{'=' * 60}
🚀 回家后操作步骤
{'=' * 60}

1. 启动后端:  bash build.sh backend --no-lint --no-test && ./ezbookkeeping server run

2. 登录系统，创建投资池账户（在现有账户体系中添加）:
   - 父账户: 投资账户 (Type=MultiSubAccounts, Category=INVESTMENT)
   - 子账户: 工资策略池 (作为投资池)
   - 记下子账户的 ID

3. 编辑此脚本 CONFIG 区域:
   - account_id: 填入投资池子账户 ID
   - token:       填入 JWT Token

4. 预览:  python convert_investment_data.py --dry-run

5. 导入:  python convert_investment_data.py --import

⚠ 已知限制:
   - quantity=0, price=0: CSV 中无确认份额/净值，需后续手动补充或
     行情 Provider 自动拉取后计算
   - conversion_out/in 的 relatedTransactionId 需创建后手动关联
   - "国泰瑞悦3个月持有期债券(FOF)" 无代码，已跳过，需手动补充 code
""")


def generate_import_guide(path, assets, transactions, conv_pairs, skipped):
    """生成 markdown 导入指引"""
    with open(path, "w", encoding="utf-8") as f:
        f.write("# 投资数据导入指引\n\n")
        f.write(f"> 自动生成，共 {len(assets)} 个资产，{len(transactions)} 条交易\n\n")

        f.write("## 1. 资产清单\n\n")
        f.write("| # | Code | Name | Type | Market | Currency |\n")
        f.write("|---|------|------|------|--------|----------|\n")
        for idx, a in enumerate(assets, 1):
            type_names = {1: "Fund", 2: "Stock", 3: "ETF", 4: "Bond", 5: "Crypto"}
            f.write(f"| {idx} | {a['code']} | {a['name']} | {type_names.get(a['type'],'?')} | CN | CNY |\n")

        f.write("\n## 2. 交易统计\n\n")
        type_counts = defaultdict(int)
        for tx in transactions:
            type_counts[tx["type"]] += 1
        f.write("| 类型 | 数量 |\n")
        f.write("|------|------|\n")
        for t in sorted(type_counts):
            f.write(f"| {TYPE_LABELS_CN.get(t, '?')} | {type_counts[t]} |\n")
        f.write(f"| **总计** | **{len(transactions)}** |\n")

        f.write(f"\n## 3. 转换配对（需手动关联 relatedTransactionId）\n\n")
        f.write("| # | 转出 | 转入 | 金额 |\n")
        f.write("|---|------|------|------|\n")
        for idx, (out_i, in_i) in enumerate(conv_pairs, 1):
            out_tx = transactions[out_i]
            in_tx = transactions[in_i]
            f.write(f"| {idx} | {out_tx.get('asset_code','?')} | {in_tx.get('asset_code','?')} "
                    f"| {out_tx['amount'] / 10000:.2f} |\n")

        if skipped:
            f.write(f"\n## 4. 跳过的记录（{len(skipped)} 条）\n\n")
            f.write("| CSV行 | 原因 |\n")
            f.write("|-------|------|\n")
            for line_no, reason in skipped:
                f.write(f"| {line_no} | {reason} |\n")


def do_import(assets, transactions, conv_pairs):
    """实际调用 API 导入数据"""
    token = CONFIG["token"]
    base = CONFIG["api_base"]

    if "REPLACE" in str(token):
        print("❌ 请先在 CONFIG 中填入有效的 JWT Token")
        sys.exit(1)
    if "REPLACE" in str(CONFIG["account_id"]):
        print("❌ 请先在 CONFIG 中填入投资池 AccountId")
        sys.exit(1)

    print(f"\n{'=' * 60}")
    print("🔄 开始导入...")
    print(f"{'=' * 60}")

    # Step 1: 创建资产，建立 code → id 映射
    code_to_id = {}
    print(f"\n--- Step 1: 创建 {len(assets)} 个资产 ---")
    for a in assets:
        payload = {
            "type": a["type"],
            "market": a["market"],
            "code": a["code"],
            "name": a["name"],
            "currency": a["currency"],
            "comment": a.get("comment", ""),
        }
        print(f"  创建: [{a['code']}] {a['name']}...", end=" ")
        result = post_json(None, f"{base}/investment/assets/add.json", payload, token)
        # 响应格式: {"result": {"id": "..."}, "success": true} 或 {"result": {"id": "...", "data": {...}}, "success": true}
        asset_id = None
        if result and result.get("success"):
            # 尝试从 result.result.id 或 result.data.id 获取
            inner = result.get("result", {})
            asset_id = inner.get("id") or inner.get("data", {}).get("id")
        if asset_id:
            asset_id = str(asset_id)
            code_to_id[a["code"]] = asset_id
            print(f"✅ id={asset_id}")
        else:
            print(f"❌ 失败")
            print(f"     响应: {json.dumps(result, ensure_ascii=False)[:200] if result else 'None'}")

    print(f"\n  资产映射: {json.dumps(code_to_id, indent=2)}")

    # Step 2: 创建交易
    # 先处理非配对交易，再处理配对交易
    print(f"\n--- Step 2: 创建 {len(transactions)} 条交易 ---")

    tx_id_map = {}  # index → created id

    # 先创建所有非 conversion 的交易
    pair_indices = set()
    for out_i, in_i in conv_pairs:
        pair_indices.add(out_i)
        pair_indices.add(in_i)

    for idx, tx in enumerate(transactions):
        if idx in pair_indices:
            continue  # 配对交易后面单独处理
        _create_transaction(idx, tx, code_to_id, base, token, tx_id_map)

    # 创建配对交易（先 out，拿到 id，再 in 并关联）
    print(f"\n  --- 配对交易 ({len(conv_pairs)} 对) ---")
    for out_i, in_i in conv_pairs:
        out_tx = transactions[out_i]
        in_tx = transactions[in_i]

        # 先创建 conversion_out
        out_id = _create_transaction(out_i, out_tx, code_to_id, base, token, tx_id_map)
        # 创建 conversion_in，关联 relatedTransactionId
        if out_id:
            in_tx_copy = dict(in_tx)
            in_tx_copy["relatedTransactionId"] = out_id
            _create_transaction(in_i, in_tx_copy, code_to_id, base, token, tx_id_map)

    print(f"\n{'=' * 60}")
    print("✅ 导入完成")
    print(f"{'=' * 60}")


def _create_transaction(idx, tx, code_to_id, base, token, tx_id_map):
    """创建单条交易记录，返回创建的 ID"""
    code = tx.get("asset_code")
    asset_id = code_to_id.get(code)
    if not asset_id:
        print(f"  ⚠ 跳过交易 {idx}: 资产 {code} 未创建")
        return None

    payload = {
        "assetId": asset_id,
        "accountId": tx["accountId"],
        "type": tx["type"],
        "tradeTime": tx["tradeTime"],
        "confirmTime": tx["confirmTime"],
        "quantity": tx["quantity"],
        "price": tx["price"],
        "amount": tx["amount"],
        "fee": tx["fee"],
        "utcOffset": tx["utcOffset"],
        "comment": tx.get("comment", ""),
    }
    if tx.get("relatedTransactionId"):
        payload["relatedTransactionId"] = tx["relatedTransactionId"]

    type_cn = TYPE_LABELS_CN.get(tx["type"], "?")
    amt = tx["amount"] / 10000
    print(f"  [{idx + 1:3d}] {type_cn:6s} {code} {amt:>10,.2f} CNY ...", end=" ")

    result = post_json(None, f"{base}/investment/transactions/add.json", payload, token)
    created_id = None
    if result and result.get("success"):
        inner = result.get("result", {})
        created_id = inner.get("id") or inner.get("data", {}).get("id")
    if created_id:
        created_id = str(created_id)
        tx_id_map[idx] = created_id
        print(f"✅ id={created_id}")
        return created_id
    else:
        print(f"❌ 失败")
        return None


if __name__ == "__main__":
    main()
