#coding=utf8

import time
from selenium import webdriver
 
driver = webdriver.Chrome()
driver.get('http://mf.91yunma.cn')

elem = driver.find_element_by_id("username")
elem.send_keys("xxxxx")
elem = driver.find_element_by_id("password")
elem.send_keys("xxxx")
elem = driver.find_element_by_id("vcode")
elem.click()
# elem = driver.find_element_by_id("btnLogin")
# elem.click()

time.sleep(8)
driver.get('http://mf.91yunma.cn/admin/qpay/get_tasks')
print driver.title
# assert "待处理订单".decode().encode('utf-8') in driver.title

isAuto = False
for i in range(100):
    driver.execute_script("getTask(6)")
    time.sleep(1)
    if isAuto == False:
        elem = driver.find_element_by_name("auto")
        elem.click()
        isAuto = True
    elem = driver.find_element_by_id("submit")
    elem.click()
    time.sleep(3)