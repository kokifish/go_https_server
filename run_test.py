import requests
import os

url = "https://localhost/uploadreport"
f = open(os.path.join("test", "file_a"), "rb")
f_html = open(os.path.join("test", "file.html"), "rb")
response = requests.post(url, files={"report_file": f, "report_html": f_html}, verify=False)  # verify="server.crt")
f.close()
f_html.close()
print(response.text)


def test_upload_v8(fname: str, tag: str, version: str):
    url = "https://localhost/uploadV8Zip"
    with open(os.path.join("test", fname), "rb") as f:
        files = {"V8zip": f}
        data = {"tag": tag, "version": version}
        response = requests.post(url, files=files, data=data, verify=False)

    print(response.text)


test_upload_v8("test_d8_zip.zip", "TEST", "122")
test_upload_v8("test_d8_2.zip", "TEST", "122")
test_upload_v8("test_d8_zip.zip", "TEST", "122")
