# holidays-jp

Simple API that provides public holidays information in Japan.

## Synopsis

### List holidays in a year

`GET /{year}` lists holidays in a year.

Example: list holidays in 2021.

```
curl https://holidays-jp.shogo82148.com/2021 | jq .
{
  "holidays": [
    {
      "date": "2021-01-01",
      "name": "元日"
    },
    {
      "date": "2021-01-11",
      "name": "成人の日"
    },
    {
      "date": "2021-02-11",
      "name": "建国記念の日"
    },
(snip)
    {
      "date": "2021-11-23",
      "name": "勤労感謝の日"
    }
  ]
}
```

### List holidays in a month

`GET /{year}/{month}` lists holidays in a year.

Example: list holidays in January 2021.

```
curl https://holidays-jp.shogo82148.com/2021/01 | jq .
{
  "holidays": [
    {
      "date": "2021-01-01",
      "name": "元日"
    },
    {
      "date": "2021-01-11",
      "name": "成人の日"
    }
  ]
}
```

### Check whether the day is a holiday

`GET /{year}/{month}/{day}` returns whether the day is a holiday.

Example: January 1st, 2021 was a holiday. The api returns information about the holiday.

```
curl https://holidays-jp.shogo82148.com/2021/01/01 | jq .
{
  "holidays": [
    {
      "date": "2021-01-01",
      "name": "元日"
    }
  ]
}
```

Example: February 1st, 2021 was not a holiday. The api returns an empty list in this case.

```
curl https://holidays-jp.shogo82148.com/2021/02/01 | jq .
{
  "holidays": []
}
```

## Data Sources

- [国民の祝日について - 内閣府](https://www8.cao.go.jp/chosei/shukujitsu/gaiyou.html) (Kokumin no Shukujitsu ni Tsuite: About Holidays in Japan - Cabinet Office, Government of Japan)

## References

- [国民の祝日に関する法律 - e-Gov 法令検索](https://elaws.e-gov.go.jp/document?lawid=323AC1000000178) (Kokumin no Shukujitsu ni kansuru Horitsu: The Law about Holidays in Japan)
- 長沢 工(1999) "日の出・日の入りの計算 天体の出没時刻の求め方" 株式会社地人書館
