from locust import HttpUser, between, task

import locust_test


class QuickstartUser(HttpUser):
    wait_time = between(3, 5)


    #登陆请求
    @task(8)
    def Login(self):

        self.client.get("/login?userid=HrKcLApw")


    #用户领取礼品
    @task(2)
    def client(self):
        data={
            "user": "HrKcLApw",
            "code": "1ZAZB42T"
        }

        self.client.post("/Client_Verify_GiftCode",json=data)

