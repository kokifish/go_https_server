import requests
import os

url = "https://localhost/uploadreport"

f = open(os.path.join("test", "file_a"), "rb")
f_html = open(os.path.join("test", "file.html"), "rb")
r = requests.post(url, files={"report_file": f, "report_html": f_html}, verify=False)  # verify="server.crt")
f.close()
f_html.close()
print(r.text)
