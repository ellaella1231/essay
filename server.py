import os
import json
import base64
import requests
from http.server import HTTPServer, SimpleHTTPRequestHandler

API_KEY = "sk-9c8474af95444cf78a1828ac228582de"
BASE_URL = "https://dashscope.aliyuncs.com/compatible-mode"
MODEL = "qwen-vl-max"

def grade_essay(base64_image, prompt_text):
    rubric_path = "backend/rubrics/Grade 7.txt"
    with open(rubric_path, "r") as f:
        rubric_text = f.read()

    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {API_KEY}"
    }
    url = f"{BASE_URL}/v1/chat/completions"

    # Step 1: OCR
    ocr_payload = {
        "model": MODEL,
        "messages": [
            {
                "role": "user",
                "content": [
                    {"type": "text", "text": "Transcribe the handwritten English text in this image exactly as written. Do not correct any mistakes."},
                    {"type": "image_url", "image_url": {"url": f"data:image/jpeg;base64,{base64_image}"}}
                ]
            }
        ]
    }

    try:
        ocr_resp = requests.post(url, headers=headers, json=ocr_payload)
        ocr_resp.raise_for_status()
        student_content = ocr_resp.json()["choices"][0]["message"]["content"]
    except Exception as e:
        return {"error": f"OCR failed: {str(e)}"}

    # Step 2: Grade
    system_prompt = f"""You are an expert English teacher. 
Grade the following student essay based on this rubric:
{rubric_text}

Essay Prompt: {prompt_text}

Student Essay:
{student_content}

You must return a JSON object strictly matching this structure:
{{
  "score": integer,
  "perfect_essay": string,
  "errors": [
    {{
      "original_segment": string,
      "suggested_segment": string,
      "explanation": string
    }}
  ]
}}"""

    grading_payload = {
        "model": MODEL,
        "response_format": {"type": "json_object"},
        "messages": [
            {"role": "system", "content": "You are an expert English teacher. You must return a JSON object strictly matching the required structure."},
            {"role": "user", "content": system_prompt}
        ]
    }

    try:
        grade_resp = requests.post(url, headers=headers, json=grading_payload)
        grade_resp.raise_for_status()
        grading_result = json.loads(grade_resp.json()["choices"][0]["message"]["content"])
        
        # Add original text for frontend display
        grading_result["original_text"] = student_content
        return grading_result
    except Exception as e:
        return {"error": f"Grading failed: {str(e)}"}

class DemoHandler(SimpleHTTPRequestHandler):
    def do_GET(self):
        if self.path == '/':
            self.path = '/demo/index.html'
        return super().do_GET()

    def do_POST(self):
        if self.path == '/api/grade':
            content_length = int(self.headers['Content-Length'])
            post_data = self.rfile.read(content_length)
            
            try:
                data = json.loads(post_data)
                base64_image = data.get('image')
                prompt_text = data.get('prompt', '')

                if not base64_image:
                    self.send_response(400)
                    self.end_headers()
                    self.wfile.write(b'{"error": "Missing image"}')
                    return

                # Remove data uri prefix if exists
                if "," in base64_image:
                    base64_image = base64_image.split(",")[1]

                result = grade_essay(base64_image, prompt_text)
                
                self.send_response(200)
                self.send_header('Content-Type', 'application/json')
                self.end_headers()
                self.wfile.write(json.dumps(result).encode('utf-8'))
                
            except Exception as e:
                self.send_response(500)
                self.end_headers()
                self.wfile.write(json.dumps({"error": str(e)}).encode('utf-8'))
        else:
            self.send_response(404)
            self.end_headers()

if __name__ == '__main__':
    server_address = ('', 8000)
    httpd = HTTPServer(server_address, DemoHandler)
    print("Serving demo at http://localhost:8000")
    httpd.serve_forever()
