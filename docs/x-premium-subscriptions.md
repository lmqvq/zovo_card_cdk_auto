# X Premium CDK

后台卡台CDK页面增加ChatGPT/X Premium产品选择，X目录只读取卡台`plans?product=x`返回的已开放套餐。

套餐：`basic_monthly`, `basic_yearly`, `premium_monthly`, `premium_yearly`, `premium_plus_monthly`, `premium_plus_yearly`。X码发码固定JP/JPY；实际支付价在兑换账号预检时获取。

用户填入含auth_token、ct0、billing_email的X Cookie JSON，不使用ChatGPT Session或邮箱密码。单笔与批量凭据解析均支持X。X Cookie不写入本站的ChatGPT账单Session缓存。

当前只支持新订阅；已有X订阅不在这里升级，自动续费由用户在X管理。依赖卡台/ACC配套版本及开关，未开放时不回落成GPT套餐。

本次未部署。
