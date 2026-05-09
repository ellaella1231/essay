import os
import base64
import json
import requests
import argparse

def process_essay(image_path, prompt_text, rubric_text, api_key, base_url, model):
    print("1. Reading image...")
    with open(image_path, "rb") as image_file:
        base64_image = base64.b64encode(image_file.read()).decode('utf-8')
    
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {api_key}"
    }

    url = f"{base_url.rstrip('/')}/v1/chat/completions"

    print("2. Calling Vision Model for Transcription...")
    ocr_payload = {
        "model": model,
        "messages": [
            {
                "role": "user",
                "content": [
                    {"type": "text", "text": "Transcribe the handwritten English text in this image exactly as written. Do not correct any mistakes."},
                    {
                        "type": "image_url",
                        "image_url": {
                            "url": f"data:image/jpeg;base64,{base64_image}"
                        }
                    }
                ]
            }
        ]
    }

    response = requests.post(url, headers=headers, json=ocr_payload)
    if response.status_code != 200:
        print("OCR Failed:", response.text)
        return

    student_content = response.json()["choices"][0]["message"]["content"]
    print("\n--- OCR Result ---")
    print(student_content)
    print("------------------\n")

    print("3. Calling Text Model for Grading and Correction...")
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
        "model": model,
        "response_format": { "type": "json_object" },
        "messages": [
            {
                "role": "system",
                "content": "You are an expert English teacher. You must return a JSON object strictly matching the required structure."
            },
            {
                "role": "user",
                "content": system_prompt
            }
        ]
    }

    response = requests.post(url, headers=headers, json=grading_payload)
    if response.status_code != 200:
        print("Grading Failed:", response.text)
        return

    grading_result = response.json()["choices"][0]["message"]["content"]
    print("\n--- Final Grading Result (JSON) ---")
    print(grading_result)
    print("-----------------------------------\n")

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--image", required=True, help="Path to essay image")
    parser.add_argument("--prompt", default="Describe your favorite hobby.", help="Essay prompt text")
    parser.add_argument("--key", required=True, help="OpenAI API Key")
    parser.add_argument("--base", default="https://api.openai.com", help="API Base URL")
    parser.add_argument("--model", default="gpt-4o", help="Model name")
    
    args = parser.parse_args()

    rubric_path = "backend/rubrics/Grade 7.txt"
    with open(rubric_path, "r") as f:
        rubric_text = f.read()
    
    process_essay(args.image, args.prompt, rubric_text, args.key, args.base, args.model)
