我是光年实验室高级招聘经理。
我在github上访问了你的开源项目，你的代码超赞。你最近有没有在看工作机会，我们在招软件开发工程师，拉钩和BOSS等招聘网站也发布了相关岗位，有公司和职位的详细信息。
我们公司在杭州，业务主要做流量增长，是很多大型互联网公司的流量顾问。公司弹性工作制，福利齐全，发展潜力大，良好的办公环境和学习氛围。
公司官网是http://www.gnlab.com,公司地址是杭州市西湖区古墩路紫金广场B座，若你感兴趣，欢迎与我联系，
电话是0571-88839161，手机号：18668131388，微信号：echo 'bGhsaGxoMTEyNAo='|base64 -D ,静待佳音。如有打扰，还请见谅，祝生活愉快工作顺利。

# locationengine

This Go library helps you to connect to Kontakt.io's MQTT broker and receive 
realtime location data of items near your Kontakt.io receiers.

After parsing the location data, the library triggers your callback functions,
notifying your application of new items appearing, items disappearing, as well
as changes in RSSI and proximity status.

Usage of the library is demonstrated in the simple [printevents](cmd/printevents) app.
You will need a Kontakt.io API key as well as the UUID of your place or receiver in
order to try this app out.

More information on the API used can be found on the 
[Kontakt.io Location Engine Monitoring page](https://developer.kontakt.io/rest-api/api-guides/location-engine-monitoring/#mqtt)
