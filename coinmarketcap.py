from selenium import webdriver
from selenium.webdriver.common.by import By 
from selenium.webdriver.common.keys import Keys
from selenium.webdriver import ActionChains
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
import time,os,random,pandas

url = "https://coinmarketcap.com/zh/?page={page_cnt}"
headers = {'User-Agent': 'User-Agent:Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36'}

browser = webdriver.Firefox()
fp = open("coins.txt", 'a', errors='ignore')

for page_cnt in range(1, 68):

    browser.get(url.format(page_cnt=page_cnt))

    browser.execute_script("window.scrollBy(0,1000)")
    time.sleep(0.5)
    browser.execute_script("window.scrollBy(0,1000)")
    time.sleep(0.5)
    browser.execute_script("window.scrollBy(0,1000)")
    time.sleep(0.5)
    browser.execute_script("window.scrollBy(0,1000)")
    time.sleep(0.5)
    browser.execute_script("window.scrollBy(0,1000)")
    time.sleep(0.5)
    browser.execute_script("window.scrollBy(0,1000)")
    time.sleep(0.5)
    browser.execute_script("window.scrollBy(0,1000)")
    time.sleep(0.5)
    browser.execute_script("window.scrollBy(0,1000)")
    time.sleep(0.5)

    df = pandas.read_html(browser.page_source)
    for i in range(100):
        try:
            rank = str(int(df[0].loc[i, '#']))
            name = df[0].loc[i, '名称'].split(rank)[0]
            print(rank, name)
            print(name, file=fp)
        except(ValueError):
            print("NaN")

            # print("Write Coin Name Fail!")
            continue

fp.close()