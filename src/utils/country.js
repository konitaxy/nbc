
const countries =[
    {
        "code": "AF",
        "label": "Afghanistan",
        "dialCode": 93,
        "phoneFormat": "070 123 4567",
        "CNName": "阿富汗"
    },
    {
        "code": "AL",
        "label": "Albania",
        "dialCode": 355,
        "phoneFormat": "066 123 4567",
        "CNName": "阿尔巴尼亚"
    },
    {
        "code": "DZ",
        "label": "Algeria",
        "dialCode": 213,
        "phoneFormat": "0551 23 45 67",
        "CNName": "阿尔及利亚"
    },
    {
        "code": "AS",
        "label": "American Samoa",
        "dialCode": 1684,
        "phoneFormat": "(684) 733-1234",
        "CNName": "美属萨摩亚"
    },
    {
        "code": "AD",
        "label": "Andorra",
        "dialCode": 376,
        "phoneFormat": "312 345",
        "CNName": "安道尔"
    },
    {
        "code": "AO",
        "label": "Angola",
        "dialCode": 244,
        "phoneFormat": "923 123 456",
        "CNName": "安哥拉"
    },
    {
        "code": "AI",
        "label": "Anguilla",
        "dialCode": 1264,
        "phoneFormat": "(264) 235-1234",
        "CNName": "安圭拉"
    },
    {
        "code": "AG",
        "label": "Antigua and Barbuda",
        "dialCode": 1268,
        "phoneFormat": "(268) 464-1234",
        "CNName": "安提瓜和巴布达"
    },
    {
        "code": "AR",
        "label": "Argentina",
        "dialCode": 54,
        "phoneFormat": "011 15-2345-6789",
        "CNName": "阿根廷"
    },
    {
        "code": "AM",
        "label": "Armenia",
        "dialCode": 374,
        "phoneFormat": "077 123456",
        "CNName": "亚美尼亚"
    },
    {
        "code": "AW",
        "label": "Aruba",
        "dialCode": 297,
        "phoneFormat": "560 1234",
        "CNName": "阿鲁巴"
    },
    {
        "code": "AU",
        "label": "Australia",
        "dialCode": 61,
        "phoneFormat": "0412 345 678",
        "CNName": "澳大利亚"
    },
    {
        "code": "AT",
        "label": "Austria",
        "dialCode": 43,
        "phoneFormat": "0664 123456",
        "CNName": "奥地利"
    },
    {
        "code": "AZ",
        "label": "Azerbaijan",
        "dialCode": 994,
        "phoneFormat": "040 123 45 67",
        "CNName": "阿塞拜疆"
    },
    {
        "code": "BH",
        "label": "Bahrain",
        "dialCode": 973,
        "phoneFormat": "3600 1234",
        "CNName": "巴林"
    },
    {
        "code": "BD",
        "label": "Bangladesh",
        "dialCode": 880,
        "phoneFormat": "01812-345678",
        "CNName": "孟加拉国"
    },
    {
        "code": "BB",
        "label": "Barbados",
        "dialCode": 1246,
        "phoneFormat": "(246) 250-1234",
        "CNName": "巴巴多斯"
    },
    {
        "code": "BY",
        "label": "Belarus",
        "dialCode": 375,
        "phoneFormat": "8 029 491-19-11",
        "CNName": "白俄罗斯"
    },
    {
        "code": "BE",
        "label": "Belgium",
        "dialCode": 32,
        "phoneFormat": "0470 12 34 56",
        "CNName": "比利时"
    },
    {
        "code": "BZ",
        "label": "Belize",
        "dialCode": 501,
        "phoneFormat": "622-1234",
        "CNName": "伯利兹"
    },
    {
        "code": "BJ",
        "label": "Benin",
        "dialCode": 229,
        "phoneFormat": "90 01 12 34",
        "CNName": "贝宁"
    },
    {
        "code": "BM",
        "label": "Bermuda",
        "dialCode": 1441,
        "phoneFormat": "(441) 370-1234",
        "CNName": "百慕大"
    },
    {
        "code": "BT",
        "label": "Bhutan",
        "dialCode": 975,
        "phoneFormat": "17 12 34 56",
        "CNName": "不丹"
    },
    {
        "code": "BO",
        "label": "Bolivia",
        "dialCode": 591,
        "phoneFormat": "71234567",
        "CNName": "玻利维亚"
    },
    {
        "code": "BA",
        "label": "Bosnia and Herzegovina",
        "dialCode": 387,
        "phoneFormat": "061 123 456",
        "CNName": "波黑"
    },
    {
        "code": "BW",
        "label": "Botswana",
        "dialCode": 267,
        "phoneFormat": "71 123 456",
        "CNName": "博茨瓦纳"
    },
    {
        "code": "BR",
        "label": "Brazil",
        "dialCode": 55,
        "phoneFormat": "(11) 96123-4567",
        "CNName": "巴西"
    },
    {
        "code": "IO",
        "label": "British Indian Ocean Territory",
        "dialCode": 246,
        "phoneFormat": "380 1234",
        "CNName": "英属印度洋领地"
    },
    {
        "code": "BN",
        "label": "Brunei Darussalam",
        "dialCode": 673,
        "phoneFormat": "712 3456",
        "CNName": "文莱"
    },
    {
        "code": "BG",
        "label": "Bulgaria",
        "dialCode": 359,
        "phoneFormat": "048 123 456",
        "CNName": "保加利亚"
    },
    {
        "code": "BF",
        "label": "Burkina Faso",
        "dialCode": 226,
        "phoneFormat": "70 12 34 56",
        "CNName": "布基纳法索"
    },
    {
        "code": "BI",
        "label": "Burundi",
        "dialCode": 257,
        "phoneFormat": "79 56 12 34",
        "CNName": "布隆迪"
    },
    {
        "code": "KH",
        "label": "Cambodia",
        "dialCode": 855,
        "phoneFormat": "091 234 567",
        "CNName": "柬埔寨"
    },
    {
        "code": "CM",
        "label": "Cameroon",
        "dialCode": 237,
        "phoneFormat": "6 71 23 45 67",
        "CNName": "喀麦隆"
    },
    {
        "code": "CA",
        "label": "Canada",
        "dialCode": 1,
        "phoneFormat": "(204) 234-5678",
        "CNName": "加拿大"
    },
    {
        "code": "CV",
        "label": "Cape Verde",
        "dialCode": 238,
        "phoneFormat": "991 12 34",
        "CNName": "佛得角"
    },
    {
        "code": "KY",
        "label": "Cayman Islands",
        "dialCode": 1345,
        "phoneFormat": "(345) 323-1234",
        "CNName": "开曼群岛"
    },
    {
        "code": "CF",
        "label": "Central African Republic",
        "dialCode": 236,
        "phoneFormat": "70 01 23 45",
        "CNName": "中非"
    },
    {
        "code": "TD",
        "label": "Chad",
        "dialCode": 235,
        "phoneFormat": "63 01 23 45",
        "CNName": "乍得"
    },
    {
        "code": "CL",
        "label": "Chile",
        "dialCode": 56,
        "phoneFormat": "09 6123 4567",
        "CNName": "智利"
    },
    {
        "code": "CN",
        "label": "China",
        "dialCode": 86,
        "phoneFormat": "131 2345 6789",
        "CNName": "中国"
    },
    {
        "code": "CO",
        "label": "Colombia",
        "dialCode": 57,
        "phoneFormat": "321 1234567",
        "CNName": "哥伦比亚"
    },
    {
        "code": "KM",
        "label": "Comoros",
        "dialCode": 269,
        "phoneFormat": "321 23 45",
        "CNName": "科摩罗"
    },
    {
        "code": "CD",
        "label": "Congo",
        "dialCode": 243,
        "phoneFormat": "0991 234 567",
        "CNName": "刚果（布）"
    },
    {
        "code": "CG",
        "label": "Congo",
        "dialCode": 242,
        "phoneFormat": "06 123 4567",
        "CNName": "刚果（金）"
    },
    {
        "code": "CK",
        "label": "Cook Islands",
        "dialCode": 682,
        "phoneFormat": "71 234",
        "CNName": "库克群岛"
    },
    {
        "code": "CR",
        "label": "Costa Rica",
        "dialCode": 506,
        "phoneFormat": "8312 3456",
        "CNName": "哥斯达黎加"
    },
    {
        "code": "CI",
        "label": "Côte d’Ivoire",
        "dialCode": 225,
        "phoneFormat": "01 23 45 67",
        "CNName": "科特迪瓦"
    },
    {
        "code": "HR",
        "label": "Croatia",
        "dialCode": 385,
        "phoneFormat": "091 234 5678",
        "CNName": "克罗地亚"
    },
    {
        "code": "CU",
        "label": "Cuba",
        "dialCode": 53,
        "phoneFormat": "05 1234567",
        "CNName": "古巴"
    },
    {
        "code": "CW",
        "label": "Curaçao",
        "dialCode": 599,
        "phoneFormat": "9 518 1234",
        "CNName": "库拉索"
    },
    {
        "code": "CY",
        "label": "Cyprus",
        "dialCode": 357,
        "phoneFormat": "96 123456",
        "CNName": "塞浦路斯"
    },
    {
        "code": "CZ",
        "label": "Czech Republic",
        "dialCode": 420,
        "phoneFormat": "601 123 456",
        "CNName": "捷克"
    },
    {
        "code": "DK",
        "label": "Denmark",
        "dialCode": 45,
        "phoneFormat": "20 12 34 56",
        "CNName": "丹麦"
    },
    {
        "code": "DJ",
        "label": "Djibouti",
        "dialCode": 253,
        "phoneFormat": "77 83 10 01",
        "CNName": "吉布提"
    },
    {
        "code": "DM",
        "label": "Dominica",
        "dialCode": 1767,
        "phoneFormat": "(767) 225-1234",
        "CNName": "多米尼克"
    },
    {
        "code": "DO",
        "label": "Dominican Republic",
        "dialCode": 1,
        "phoneFormat": "(809) 234-5678",
        "CNName": "多米尼加"
    },
    {
        "code": "EC",
        "label": "Ecuador",
        "dialCode": 593,
        "phoneFormat": "099 123 4567",
        "CNName": "厄瓜多尔"
    },
    {
        "code": "EG",
        "label": "Egypt",
        "dialCode": 20,
        "phoneFormat": "0100 123 4567",
        "CNName": "埃及"
    },
    {
        "code": "SV",
        "label": "El Salvador",
        "dialCode": 503,
        "phoneFormat": "7012 3456",
        "CNName": "萨尔瓦多"
    },
    {
        "code": "GQ",
        "label": "Equatorial Guinea",
        "dialCode": 240,
        "phoneFormat": "222 123 456",
        "CNName": "赤道几内亚"
    },
    {
        "code": "ER",
        "label": "Eritrea",
        "dialCode": 291,
        "phoneFormat": "07 123 456",
        "CNName": "厄立特里亚"
    },
    {
        "code": "EE",
        "label": "Estonia",
        "dialCode": 372,
        "phoneFormat": "5123 4567",
        "CNName": "爱沙尼亚"
    },
    {
        "code": "ET",
        "label": "Ethiopia",
        "dialCode": 251,
        "phoneFormat": "091 123 4567",
        "CNName": "埃塞俄比亚"
    },
    {
        "code": "FK",
        "label": "Falkland Islands",
        "dialCode": 500,
        "phoneFormat": "51234",
        "CNName": "福克兰群岛（马尔维纳斯）"
    },
    {
        "code": "FO",
        "label": "Faroe Islands",
        "dialCode": 298,
        "phoneFormat": "211234",
        "CNName": "法罗群岛"
    },
    {
        "code": "FJ",
        "label": "Fiji",
        "dialCode": 679,
        "phoneFormat": "701 2345",
        "CNName": "斐济"
    },
    {
        "code": "FI",
        "label": "Finland",
        "dialCode": 358,
        "phoneFormat": "041 2345678",
        "CNName": "芬兰"
    },
    {
        "code": "FR",
        "label": "France",
        "dialCode": 33,
        "phoneFormat": "06 12 34 56 78",
        "CNName": "法国"
    },
    {
        "code": "GF",
        "label": "French Guiana",
        "dialCode": 594,
        "phoneFormat": "0694 20 12 34",
        "CNName": "法属圭亚那"
    },
    {
        "code": "PF",
        "label": "French Polynesia",
        "dialCode": 689,
        "phoneFormat": "87 12 34 56",
        "CNName": "法属波利尼西亚"
    },
    {
        "code": "GA",
        "label": "Gabon",
        "dialCode": 241,
        "phoneFormat": "06 03 12 34",
        "CNName": "加蓬"
    },
    {
        "code": "GE",
        "label": "Georgia",
        "dialCode": 995,
        "phoneFormat": "555 12 34 56",
        "CNName": "格鲁吉亚"
    },
    {
        "code": "DE",
        "label": "Germany",
        "dialCode": 49,
        "phoneFormat": "01512 3456789",
        "CNName": "德国"
    },
    {
        "code": "GH",
        "label": "Ghana",
        "dialCode": 233,
        "phoneFormat": "023 123 4567",
        "CNName": "加纳"
    },
    {
        "code": "GI",
        "label": "Gibraltar",
        "dialCode": 350,
        "phoneFormat": "57123456",
        "CNName": "直布罗陀"
    },
    {
        "code": "GR",
        "label": "Greece",
        "dialCode": 30,
        "phoneFormat": "691 234 5678",
        "CNName": "希腊"
    },
    {
        "code": "GL",
        "label": "Greenland",
        "dialCode": 299,
        "phoneFormat": "22 12 34",
        "CNName": "格陵兰"
    },
    {
        "code": "GD",
        "label": "Grenada",
        "dialCode": 1473,
        "phoneFormat": "(473) 403-1234",
        "CNName": "格林纳达"
    },
    {
        "code": "GP",
        "label": "Guadeloupe",
        "dialCode": 590,
        "phoneFormat": "0690 30-1234",
        "CNName": "瓜德罗普"
    },
    {
        "code": "GU",
        "label": "Guam",
        "dialCode": 1671,
        "phoneFormat": "(671) 300-1234",
        "CNName": "关岛"
    },
    {
        "code": "GT",
        "label": "Guatemala",
        "dialCode": 502,
        "phoneFormat": "5123 4567",
        "CNName": "危地马拉"
    },
    {
        "code": "GG",
        "label": "Guernsey",
        "dialCode": 1481,
        "phoneFormat": "07781 123456",
        "CNName": "根西岛"
    },
    {
        "code": "GN",
        "label": "Guinea",
        "dialCode": 224,
        "phoneFormat": "601 12 34 56",
        "CNName": "几内亚"
    },
    {
        "code": "GW",
        "label": "Guinea Bissau",
        "dialCode": 245,
        "phoneFormat": "955 012 345",
        "CNName": "几内亚比绍"
    },
    {
        "code": "GY",
        "label": "Guyana",
        "dialCode": 592,
        "phoneFormat": "609 1234",
        "CNName": "圭亚那"
    },
    {
        "code": "HT",
        "label": "Haiti",
        "dialCode": 509,
        "phoneFormat": "34 10 1234",
        "CNName": "海地"
    },
    {
        "code": "HN",
        "label": "Honduras",
        "dialCode": 504,
        "phoneFormat": "9123-4567",
        "CNName": "洪都拉斯"
    },
    {
        "code": "HK",
        "label": "Hong Kong SAR",
        "dialCode": 852,
        "phoneFormat": "5123 4567",
        "CNName": "香港"
    },
    {
        "code": "HU",
        "label": "Hungary",
        "dialCode": 36,
        "phoneFormat": "(20) 123 4567",
        "CNName": "匈牙利"
    },
    {
        "code": "IS",
        "label": "Iceland",
        "dialCode": 354,
        "phoneFormat": "611 1234",
        "CNName": "冰岛"
    },
    {
        "code": "IN",
        "label": "India",
        "dialCode": 91,
        "phoneFormat": "099876 54321",
        "CNName": "印度"
    },
    {
        "code": "ID",
        "label": "Indonesia",
        "dialCode": 62,
        "phoneFormat": "0812-345-678",
        "CNName": "印度尼西亚"
    },
    {
        "code": "IR",
        "label": "Iran",
        "dialCode": 98,
        "phoneFormat": "0912 345 6789",
        "CNName": "伊朗"
    },
    {
        "code": "IQ",
        "label": "Iraq",
        "dialCode": 964,
        "phoneFormat": "0791 234 5678",
        "CNName": "伊拉克"
    },
    {
        "code": "IE",
        "label": "Ireland",
        "dialCode": 353,
        "phoneFormat": "085 012 3456",
        "CNName": "爱尔兰"
    },
    {
        "code": "IL",
        "label": "Israel",
        "dialCode": 972,
        "phoneFormat": "050-123-4567",
        "CNName": "以色列"
    },
    {
        "code": "IT",
        "label": "Italy",
        "dialCode": 39,
        "phoneFormat": "312 345 6789",
        "CNName": "意大利"
    },
    {
        "code": "JM",
        "label": "Jamaica",
        "dialCode": 1876,
        "phoneFormat": "(876) 210-1234",
        "CNName": "牙买加"
    },
    {
        "code": "JP",
        "label": "Japan",
        "dialCode": 81,
        "phoneFormat": "090-1234-5678",
        "CNName": "日本"
    },
    {
        "code": "JO",
        "label": "Jordan",
        "dialCode": 962,
        "phoneFormat": "07 9012 3456",
        "CNName": "约旦"
    },
    {
        "code": "KZ",
        "label": "Kazakhstan",
        "dialCode": 7,
        "phoneFormat": "8 (771) 000 9998",
        "CNName": "哈萨克斯坦"
    },
    {
        "code": "KE",
        "label": "Kenya",
        "dialCode": 254,
        "phoneFormat": "0712 123456",
        "CNName": "肯尼亚"
    },
    {
        "code": "KI",
        "label": "Kiribati",
        "dialCode": 686,
        "phoneFormat": "72012345",
        "CNName": "基里巴斯"
    },
    {
        "code": "XK",
        "label": "Kosovo",
        "dialCode": 383,
        "phoneFormat": "",
        "CNName": "科索沃"
    },
    {
        "code": "KW",
        "label": "Kuwait",
        "dialCode": 965,
        "phoneFormat": "500 12345",
        "CNName": "科威特"
    },
    {
        "code": "KG",
        "label": "Kyrgyzstan",
        "dialCode": 996,
        "phoneFormat": "0700 123 456",
        "CNName": "吉尔吉斯斯坦"
    },
    {
        "code": "LA",
        "label": "Laos",
        "dialCode": 856,
        "phoneFormat": "020 23 123 456",
        "CNName": "老挝"
    },
    {
        "code": "LV",
        "label": "Latvia",
        "dialCode": 371,
        "phoneFormat": "21 234 567",
        "CNName": "拉脱维亚"
    },
    {
        "code": "LB",
        "label": "Lebanon",
        "dialCode": 961,
        "phoneFormat": "71 123 456",
        "CNName": "黎巴嫩"
    },
    {
        "code": "LS",
        "label": "Lesotho",
        "dialCode": 266,
        "phoneFormat": "5012 3456",
        "CNName": "莱索托"
    },
    {
        "code": "LR",
        "label": "Liberia",
        "dialCode": 231,
        "phoneFormat": "077 012 3456",
        "CNName": "利比里亚"
    },
    {
        "code": "LY",
        "label": "Libya",
        "dialCode": 218,
        "phoneFormat": "091-2345678",
        "CNName": "利比亚"
    },
    {
        "code": "LI",
        "label": "Liechtenstein",
        "dialCode": 423,
        "phoneFormat": "660 234 567",
        "CNName": "列支敦士登"
    },
    {
        "code": "LT",
        "label": "Lithuania",
        "dialCode": 370,
        "phoneFormat": "(8-612) 34567",
        "CNName": "立陶宛"
    },
    {
        "code": "LU",
        "label": "Luxembourg",
        "dialCode": 352,
        "phoneFormat": "628 123 456",
        "CNName": "卢森堡"
    },
    {
        "code": "MO",
        "label": "Macao",
        "dialCode": 853,
        "phoneFormat": "6612 3456",
        "CNName": "澳门"
    },
    {
        "code": "MK",
        "label": "Macedonia",
        "dialCode": 389,
        "phoneFormat": "072 345 678",
        "CNName": "前南马其顿"
    },
    {
        "code": "MG",
        "label": "Madagascar",
        "dialCode": 261,
        "phoneFormat": "032 12 345 67",
        "CNName": "马达加斯加"
    },
    {
        "code": "MW",
        "label": "Malawi",
        "dialCode": 265,
        "phoneFormat": "0991 23 45 67",
        "CNName": "马拉维"
    },
    {
        "code": "MY",
        "label": "Malaysia",
        "dialCode": 60,
        "phoneFormat": "012-345 6789",
        "CNName": "马来西亚"
    },
    {
        "code": "MV",
        "label": "Maldives",
        "dialCode": 960,
        "phoneFormat": "771-2345",
        "CNName": "马尔代夫"
    },
    {
        "code": "ML",
        "label": "Mali",
        "dialCode": 223,
        "phoneFormat": "65 01 23 45",
        "CNName": "马里"
    },
    {
        "code": "MT",
        "label": "Malta",
        "dialCode": 356,
        "phoneFormat": "9696 1234",
        "CNName": "马耳他"
    },
    {
        "code": "MH",
        "label": "Marshall Islands",
        "dialCode": 692,
        "phoneFormat": "235-1234",
        "CNName": "马绍尔群岛"
    },
    {
        "code": "MQ",
        "label": "Martinique",
        "dialCode": 596,
        "phoneFormat": "0696 20 12 34",
        "CNName": "马提尼克"
    },
    {
        "code": "MR",
        "label": "Mauritania",
        "dialCode": 222,
        "phoneFormat": "22 12 34 56",
        "CNName": "毛里塔尼亚"
    },
    {
        "code": "MU",
        "label": "Mauritius",
        "dialCode": 230,
        "phoneFormat": "5251 2345",
        "CNName": "毛里求斯"
    },
    {
        "code": "YT",
        "label": "Mayotte",
        "dialCode": 262,
        "phoneFormat": "0639 12 34 56",
        "CNName": "马约特"
    },
    {
        "code": "MX",
        "label": "Mexico",
        "dialCode": 52,
        "phoneFormat": "044 222 123 4567",
        "CNName": "墨西哥"
    },
    {
        "code": "FM",
        "label": "Micronesia",
        "dialCode": 691,
        "phoneFormat": "350 1234",
        "CNName": "密克罗尼西亚"
    },
    {
        "code": "MD",
        "label": "Moldova",
        "dialCode": 373,
        "phoneFormat": "0621 12 345",
        "CNName": "摩尔多瓦"
    },
    {
        "code": "MC",
        "label": "Monaco",
        "dialCode": 377,
        "phoneFormat": "06 12 34 56 78",
        "CNName": "摩纳哥"
    },
    {
        "code": "MN",
        "label": "Mongolia",
        "dialCode": 976,
        "phoneFormat": "8812 3456",
        "CNName": "蒙古"
    },
    {
        "code": "ME",
        "label": "Montenegro",
        "dialCode": 382,
        "phoneFormat": "067 622 901",
        "CNName": "黑山"
    },
    {
        "code": "MS",
        "label": "Montserrat",
        "dialCode": 1664,
        "phoneFormat": "(664) 492-3456",
        "CNName": "蒙特塞拉特"
    },
    {
        "code": "MA",
        "label": "Morocco",
        "dialCode": 212,
        "phoneFormat": "0650-123456",
        "CNName": "摩洛哥"
    },
    {
        "code": "MZ",
        "label": "Mozambique",
        "dialCode": 258,
        "phoneFormat": "82 123 4567",
        "CNName": "莫桑比克"
    },
    {
        "code": "MM",
        "label": "Myanmar",
        "dialCode": 95,
        "phoneFormat": "09 212 3456",
        "CNName": "缅甸"
    },
    {
        "code": "NA",
        "label": "Namibia",
        "dialCode": 264,
        "phoneFormat": "081 123 4567",
        "CNName": "纳米尼亚"
    },
    {
        "code": "NR",
        "label": "Nauru",
        "dialCode": 674,
        "phoneFormat": "555 1234",
        "CNName": "瑙鲁"
    },
    {
        "code": "NP",
        "label": "Nepal",
        "dialCode": 977,
        "phoneFormat": "984-1234567",
        "CNName": "尼泊尔"
    },
    {
        "code": "NL",
        "label": "Netherlands",
        "dialCode": 31,
        "phoneFormat": "06 12345678",
        "CNName": "荷兰"
    },
    {
        "code": "NC",
        "label": "Calédonie)",
        "dialCode": 687,
        "phoneFormat": "75.12.34",
        "CNName": "新喀里多尼亚-New Caledonia (Nouvelle"
    },
    {
        "code": "NZ",
        "label": "New Zealand",
        "dialCode": 64,
        "phoneFormat": "021 123 4567",
        "CNName": "新西兰"
    },
    {
        "code": "NI",
        "label": "Nicaragua",
        "dialCode": 505,
        "phoneFormat": "8123 4567",
        "CNName": "尼加拉瓜"
    },
    {
        "code": "NE",
        "label": "Niger",
        "dialCode": 227,
        "phoneFormat": "93 12 34 56",
        "CNName": "尼日尔"
    },
    {
        "code": "NG",
        "label": "Nigeria",
        "dialCode": 234,
        "phoneFormat": "0802 123 4567",
        "CNName": "尼日利亚"
    },
    {
        "code": "NU",
        "label": "Niue",
        "dialCode": 683,
        "phoneFormat": "1234",
        "CNName": "纽埃"
    },
    {
        "code": "NF",
        "label": "Norfolk Island",
        "dialCode": 672,
        "phoneFormat": "3 81234",
        "CNName": "诺福克岛"
    },
    {
        "code": "KP",
        "label": "North Korea",
        "dialCode": 850,
        "phoneFormat": "0192 123 4567",
        "CNName": "朝鲜"
    },
    {
        "code": "MP",
        "label": "Northern Mariana Islands",
        "dialCode": 1670,
        "phoneFormat": "(670) 234-5678",
        "CNName": "北马里亚纳"
    },
    {
        "code": "NO",
        "label": "Norway",
        "dialCode": 47,
        "phoneFormat": "406 12 345",
        "CNName": "挪威"
    },
    {
        "code": "OM",
        "label": "Oman",
        "dialCode": 968,
        "phoneFormat": "9212 3456",
        "CNName": "阿曼"
    },
    {
        "code": "PK",
        "label": "Pakistan",
        "dialCode": 92,
        "phoneFormat": "0301 2345678",
        "CNName": "巴基斯坦"
    },
    {
        "code": "PW",
        "label": "Palau",
        "dialCode": 680,
        "phoneFormat": "620 1234",
        "CNName": "帕劳"
    },
    {
        "code": "PS",
        "label": "Palestine",
        "dialCode": 970,
        "phoneFormat": "0599 123 456",
        "CNName": "巴勒斯坦"
    },
    {
        "code": "PA",
        "label": "Panama",
        "dialCode": 507,
        "phoneFormat": "6001-2345",
        "CNName": "巴拿马"
    },
    {
        "code": "PG",
        "label": "Papua New Guinea",
        "dialCode": 675,
        "phoneFormat": "681 2345",
        "CNName": "巴布亚新几内亚"
    },
    {
        "code": "PY",
        "label": "Paraguay",
        "dialCode": 595,
        "phoneFormat": "0961 456789",
        "CNName": "巴拉圭"
    },
    {
        "code": "PE",
        "label": "Peru",
        "dialCode": 51,
        "phoneFormat": "912 345 678",
        "CNName": "秘鲁"
    },
    {
        "code": "PH",
        "label": "Philippines",
        "dialCode": 63,
        "phoneFormat": "0905 123 4567",
        "CNName": "菲律宾"
    },
    {
        "code": "PL",
        "label": "Poland",
        "dialCode": 48,
        "phoneFormat": "512 345 678",
        "CNName": "皮特凯恩"
    },
    {
        "code": "PT",
        "label": "Portugal",
        "dialCode": 351,
        "phoneFormat": "912 345 678",
        "CNName": "葡萄牙"
    },
    {
        "code": "PR",
        "label": "Puerto Rico",
        "dialCode": 1,
        "phoneFormat": "(787) 234-5678",
        "CNName": "波多黎各"
    },
    {
        "code": "QA",
        "label": "Qatar",
        "dialCode": 974,
        "phoneFormat": "3312 3456",
        "CNName": "卡塔尔"
    },
    {
        "code": "RE",
        "label": "Réunion",
        "dialCode": 262,
        "phoneFormat": "0692 12 34 56",
        "CNName": "留尼汪"
    },
    {
        "code": "RO",
        "label": "Romania",
        "dialCode": 40,
        "phoneFormat": "0712 345 678",
        "CNName": "罗马尼亚"
    },
    {
        "code": "RU",
        "label": "Russia",
        "dialCode": 7,
        "phoneFormat": "8 (912) 345-67-89",
        "CNName": "俄罗斯"
    },
    {
        "code": "RW",
        "label": "Rwanda",
        "dialCode": 250,
        "phoneFormat": "0720 123 456",
        "CNName": "卢旺达"
    },
    {
        "code": "SH",
        "label": "Saint Helena",
        "dialCode": 290,
        "phoneFormat": "51234",
        "CNName": "圣赫勒拿"
    },
    {
        "code": "KN",
        "label": "Saint Kitts and Nevis",
        "dialCode": 1869,
        "phoneFormat": "(869) 765-2917",
        "CNName": "圣基茨和尼维斯"
    },
    {
        "code": "LC",
        "label": "Saint Lucia",
        "dialCode": 1758,
        "phoneFormat": "(758) 284-5678",
        "CNName": "圣卢西亚"
    },
    {
        "code": "MF",
        "label": "Martin)",
        "dialCode": 590,
        "phoneFormat": "0690 30-1234",
        "CNName": "圣马丁岛-Saint Martin (Saint"
    },
    {
        "code": "PM",
        "label": "Saint Pierre and Miquelon",
        "dialCode": 508,
        "phoneFormat": "055 12 34",
        "CNName": "圣皮埃尔和密克隆"
    },
    {
        "code": "VC",
        "label": "Saint Vincent and the Grenadines",
        "dialCode": 1784,
        "phoneFormat": "(784) 430-1234",
        "CNName": "圣文森特和格林纳丁斯"
    },
    {
        "code": "WS",
        "label": "Samoa",
        "dialCode": 685,
        "phoneFormat": "601234",
        "CNName": "萨摩亚"
    },
    {
        "code": "SM",
        "label": "San Marino",
        "dialCode": 378,
        "phoneFormat": "66 66 12 12",
        "CNName": "圣马力诺"
    },
    {
        "code": "ST",
        "label": "São Tomé and Príncipe",
        "dialCode": 239,
        "phoneFormat": "981 2345",
        "CNName": "圣多美和普林西比"
    },
    {
        "code": "SA",
        "label": "Saudi Arabia",
        "dialCode": 966,
        "phoneFormat": "051 234 5678",
        "CNName": "沙特阿拉伯"
    },
    {
        "code": "SN",
        "label": "Senegal",
        "dialCode": 221,
        "phoneFormat": "70 123 45 67",
        "CNName": "塞内加尔"
    },
    {
        "code": "RS",
        "label": "Serbia",
        "dialCode": 381,
        "phoneFormat": "060 1234567",
        "CNName": "塞尔维亚"
    },
    {
        "code": "SC",
        "label": "Seychelles",
        "dialCode": 248,
        "phoneFormat": "2 510 123",
        "CNName": "塞舌尔"
    },
    {
        "code": "SL",
        "label": "Sierra Leone",
        "dialCode": 232,
        "phoneFormat": "(025) 123456",
        "CNName": "塞拉利昂"
    },
    {
        "code": "SG",
        "label": "Singapore",
        "dialCode": 65,
        "phoneFormat": "8123 4567",
        "CNName": "新加坡"
    },
    {
        "code": "SX",
        "label": "Sint Maarten",
        "dialCode": 1721,
        "phoneFormat": "(721) 520-5678",
        "CNName": "荷属圣马丁"
    },
    {
        "code": "SK",
        "label": "Slovakia",
        "dialCode": 421,
        "phoneFormat": "0912 123 456",
        "CNName": "斯洛伐克"
    },
    {
        "code": "SI",
        "label": "Slovenia",
        "dialCode": 386,
        "phoneFormat": "031 234 567",
        "CNName": "斯洛文尼亚"
    },
    {
        "code": "SB",
        "label": "Solomon Islands",
        "dialCode": 677,
        "phoneFormat": "74 21234",
        "CNName": "所罗门群岛"
    },
    {
        "code": "SO",
        "label": "Somalia",
        "dialCode": 252,
        "phoneFormat": "7 1123456",
        "CNName": "索马里"
    },
    {
        "code": "ZA",
        "label": "South Africa",
        "dialCode": 27,
        "phoneFormat": "071 123 4567",
        "CNName": "南非"
    },
    {
        "code": "KR",
        "label": "South Korea",
        "dialCode": 82,
        "phoneFormat": "010-0000-0000",
        "CNName": "韩国"
    },
    {
        "code": "SS",
        "label": "South Sudan",
        "dialCode": 211,
        "phoneFormat": "0977 123 456",
        "CNName": "南苏丹"
    },
    {
        "code": "ES",
        "label": "Spain",
        "dialCode": 34,
        "phoneFormat": "612 34 56 78",
        "CNName": "西班牙"
    },
    {
        "code": "LK",
        "label": "Sri Lanka",
        "dialCode": 94,
        "phoneFormat": "071 234 5678",
        "CNName": "斯里兰卡"
    },
    {
        "code": "SD",
        "label": "Sudan",
        "dialCode": 249,
        "phoneFormat": "091 123 1234",
        "CNName": "苏丹"
    },
    {
        "code": "SR",
        "label": "Suriname",
        "dialCode": 597,
        "phoneFormat": "741-2345",
        "CNName": "苏里南"
    },
    {
        "code": "SZ",
        "label": "Swaziland",
        "dialCode": 268,
        "phoneFormat": "7612 3456",
        "CNName": "斯威士兰"
    },
    {
        "code": "SE",
        "label": "Sweden",
        "dialCode": 46,
        "phoneFormat": "070-123 45 67",
        "CNName": "瑞典"
    },
    {
        "code": "CH",
        "label": "Switzerland",
        "dialCode": 41,
        "phoneFormat": "078 123 45 67",
        "CNName": "瑞士"
    },
    {
        "code": "SY",
        "label": "Syria",
        "dialCode": 963,
        "phoneFormat": "0944 567 890",
        "CNName": "叙利亚"
    },
    {
        "code": "TW",
        "label": "Taiwan",
        "dialCode": 886,
        "phoneFormat": "0912 345 678",
        "CNName": "台湾"
    },
    {
        "code": "TJ",
        "label": "Tajikistan",
        "dialCode": 992,
        "phoneFormat": "(8) 917 12 3456",
        "CNName": "塔吉克斯坦"
    },
    {
        "code": "TZ",
        "label": "Tanzania",
        "dialCode": 255,
        "phoneFormat": "0621 234 567",
        "CNName": "坦桑尼亚"
    },
    {
        "code": "TH",
        "label": "Thailand",
        "dialCode": 66,
        "phoneFormat": "081 234 5678",
        "CNName": "泰国"
    },
    {
        "code": "BS",
        "label": "The Bahamas",
        "dialCode": 1242,
        "phoneFormat": "(242) 359-1234",
        "CNName": "巴哈马"
    },
    {
        "code": "GM",
        "label": "The Gambia",
        "dialCode": 220,
        "phoneFormat": "301 2345",
        "CNName": "冈比亚"
    },
    {
        "code": "TL",
        "label": "Leste",
        "dialCode": 670,
        "phoneFormat": "7721 2345",
        "CNName": "东帝汶-Timor"
    },
    {
        "code": "TG",
        "label": "Togo",
        "dialCode": 228,
        "phoneFormat": "90 11 23 45",
        "CNName": "多哥"
    },
    {
        "code": "TK",
        "label": "Tokelau",
        "dialCode": 690,
        "phoneFormat": "7290",
        "CNName": "托克劳"
    },
    {
        "code": "TO",
        "label": "Tonga",
        "dialCode": 676,
        "phoneFormat": "771 5123",
        "CNName": "汤加"
    },
    {
        "code": "TT",
        "label": "Trinidad and Tobago",
        "dialCode": 1868,
        "phoneFormat": "(868) 291-1234",
        "CNName": "特立尼达和多巴哥"
    },
    {
        "code": "TN",
        "label": "Tunisia",
        "dialCode": 216,
        "phoneFormat": "20 123 456",
        "CNName": "突尼斯"
    },
    {
        "code": "TR",
        "label": "Turkey",
        "dialCode": 90,
        "phoneFormat": "0501 234 56 78",
        "CNName": "土耳其"
    },
    {
        "code": "TM",
        "label": "Turkmenistan",
        "dialCode": 993,
        "phoneFormat": "8 66 123456",
        "CNName": "土库曼斯坦"
    },
    {
        "code": "TC",
        "label": "Turks and Caicos Islands",
        "dialCode": 1649,
        "phoneFormat": "(649) 231-1234",
        "CNName": "特克斯和凯科斯群岛"
    },
    {
        "code": "TV",
        "label": "Tuvalu",
        "dialCode": 688,
        "phoneFormat": "901234",
        "CNName": "图瓦卢"
    },
    {
        "code": "US",
        "label": "United States",
        "dialCode": 1,
        "phoneFormat": "(201) 555-0123",
        "CNName": "美国"
    },
    {
        "code": "GB",
        "label": "United Kingdom",
        "dialCode": 44,
        "phoneFormat": "07400 123456",
        "CNName": "英国"
    },
    {
        "code": "UG",
        "label": "Uganda",
        "dialCode": 256,
        "phoneFormat": "0712 345678",
        "CNName": "乌干达"
    },
    {
        "code": "UA",
        "label": "Ukraine",
        "dialCode": 380,
        "phoneFormat": "039 123 4567",
        "CNName": "乌克兰"
    },
    {
        "code": "AE",
        "label": "United Arab Emirates",
        "dialCode": 971,
        "phoneFormat": "050 123 4567",
        "CNName": "阿拉伯联合酋长国"
    },
    {
        "code": "UY",
        "label": "Uruguay",
        "dialCode": 598,
        "phoneFormat": "094 231 234",
        "CNName": "乌拉圭"
    },
    {
        "code": "UZ",
        "label": "Uzbekistan",
        "dialCode": 998,
        "phoneFormat": "8 91 234 56 78",
        "CNName": "乌兹别克斯坦"
    },
    {
        "code": "VU",
        "label": "Vanuatu",
        "dialCode": 678,
        "phoneFormat": "591 2345",
        "CNName": "瓦努阿图"
    },
    {
        "code": "VA",
        "label": "Vatican City",
        "dialCode": 379,
        "phoneFormat": "312 345 6789",
        "CNName": "梵蒂冈城"
    },
    {
        "code": "VE",
        "label": "Venezuela",
        "dialCode": 58,
        "phoneFormat": "0412-1234567",
        "CNName": "委内瑞拉"
    },
    {
        "code": "VN",
        "label": "Vietnam",
        "dialCode": 84,
        "phoneFormat": "091 234 56 78",
        "CNName": "越南"
    },
    {
        "code": "WF",
        "label": "Wallis and Futuna",
        "dialCode": 681,
        "phoneFormat": "50 12 34",
        "CNName": "瓦利斯和富图纳"
    },
    {
        "code": "EH",
        "label": "Western Sahara",
        "dialCode": 212,
        "phoneFormat": "0650-123456",
        "CNName": "西撒哈拉"
    },
    {
        "code": "YE",
        "label": "Yemen",
        "dialCode": 967,
        "phoneFormat": "0712 345 678",
        "CNName": "也门"
    },
    {
        "code": "ZM",
        "label": "Zambia",
        "dialCode": 260,
        "phoneFormat": "095 5123456",
        "CNName": "赞比亚"
    },
    {
        "code": "ZW",
        "label": "Zimbabwe",
        "dialCode": 263,
        "phoneFormat": "071 123 4567",
        "CNName": "津巴布韦"
    }
]
function findByCode(code) {
  return countries.find(item => item.code === code)
}
export {countries,findByCode}