# mortis
Go 언어만 있는 KYC HANA to midi 컨버터
## 개요
모르티스는 금영 노래방기기 (노래반주기)에서 사용하는 KYC HANA를 midi 파일로 변환해주는 컨버터 입니다.

de_lzah.c 소스코드를 마이크로소프트 코파일럿의 도움을 받아 순수한 Go언어로 변환하여 더이상 CGO와 C언어 컴파일러가 필요하지 않습니다.

## 주의사항 
변환하긴 했지만 그럼에도 [QuickBMS](http://aluigi.altervista.org/quickbms.htm)의 de_lzah.c 소스코드를 기반으로 하고 있어서 여전히 GPL-2.0 license를 따릅니다.
