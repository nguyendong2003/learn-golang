import json
import os
import subprocess
import time

# --- CẤU HÌNH CHẠY QUA NPX (KHÔNG CẦN CLONE REPO) ---
# Gọi trực tiếp package từ registry của npm giống hệt cách Antigravity đang cấu hình
COMMAND = ["npx", "-y", "@vkhanhqui/figma-mcp-go"]

# # --- CẤU HÌNH ---
# # Thay thế bằng lệnh chạy figma-mcp-go của bạn (hoặc đường dẫn đến file .exe/.bin đã build)
# COMMAND = ["go", "run", "./cmd/figma-mcp-go/main.go"]


# Biến toàn cục để tối ưu việc giữ kết nối xuyên suốt, tránh mất Context dữ liệu của Figma
_mcp_process = None

def call_mcp_tool(tool_name, arguments={}):
    global _mcp_process
    
    # Khởi chạy một lần duy nhất và dùng chung cho tất cả các tool để giữ trạng thái Session
    if _mcp_process is None or _mcp_process.poll() is not None:
        _mcp_process = subprocess.Popen(
            COMMAND,
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            bufsize=1 # Line buffered
        )

    # 2. Tạo Request theo chuẩn JSON-RPC của giao thức MCP
    # Để kiểm tra công cụ, chúng ta gửi method "tools/call"
    payload = {
        "jsonrpc": "2.0",
        "id": int(time.time() * 1000), # Dùng timestamp làm ID để không trùng lặp
        "method": "tools/call",
        "params": {"name": tool_name, "arguments": arguments},
    }

    try:
        # 3. Gửi request vào stdin của server Go
        input_data = json.dumps(payload) + "\n"
        print(f"--> Đang gửi lệnh gọi tool '{tool_name}' tới Figma MCP Go...")

        _mcp_process.stdin.write(input_data)
        _mcp_process.stdin.flush()

        # Đọc một dòng phản hồi từ stdout
        stdout_data = _mcp_process.stdout.readline()

        # 4. Xử lý kết quả trả về từ stdout
        if stdout_data:
            response_json = json.loads(stdout_data.strip())
            return response_json
        else:
            print("Không nhận được phản hồi từ stdout.")
            return None

    except Exception as e:
        print(f"Có lỗi xảy ra khi gọi tool {tool_name}: {e}")
        return None

def sort_json_by_geometry(data):
    """
    Hàm đệ quy sắp xếp nâng cấp: 
    Ưu tiên so sánh tọa độ TÂM trục dọc (Center Y) để xử lý chuẩn xác các phần tử căn giữa (Align Center).
    Nếu nằm cùng một hàng ngang (Tâm sau chuẩn hóa bằng nhau), xếp từ trái sang phải theo x.
    """
    if isinstance(data, list):
        # Duyệt qua từng phần tử trong list để sắp xếp children trước
        for item in data:
            sort_json_by_geometry(item)
        
        # Tiến hành sắp xếp các phần tử cùng cấp trong list
        # Định nghĩa hàm key để lấy tọa độ an toàn, tránh trường hợp node không có bounds
        def get_sort_key(node):
            if isinstance(node, dict) and "bounds" in node:
                b = node["bounds"]
                x = b.get("x", 0)
                y = b.get("y", 0)
                height = b.get("height", 0)
                
                # 1. Tính toán TÂM của trục dọc (Center Y = y + height/2) thay vì dùng tọa độ đỉnh y.
                # Điều này giúp các phần tử lệch height (như Text 21px và Icon 16px) nhận diện đúng là đang căn giữa nhau.
                center_y = y + (height / 2)
                
                # 2. Làm tròn tâm về hệ lưới 4px để triệt tiêu hoàn toàn sai số sub-pixel do designer kéo tay.
                # Khi chia cho 4, làm tròn rồi nhân ngược lại, các tâm sát nhau (22.0 và 22.5) sẽ quy về cùng một mốc.
                normalized_center_y = round(center_y / 4) * 4
                
                # Trả về bộ khóa (Tuple): Ưu tiên xếp theo Tâm dọc trước, nếu hòa nhau sẽ xếp từ trái sang phải theo x
                return (normalized_center_y, round(x))
            return (0, 0)
        
        data.sort(key=get_sort_key)

    elif isinstance(data, dict):
        # Nếu node có children, đệ quy sâu vào trong để sắp xếp cụm con của nó
        if "children" in data and isinstance(data["children"], list):
            sort_json_by_geometry(data["children"])
            
    return data

# def sort_json_by_geometry(data):
#     """
#     Hàm đệ quy sắp xếp tất cả các phần tử (và children của chúng) 
#     theo thứ tự từ trên xuống dưới (y), nếu y bằng nhau thì xếp từ trái sang phải (x).
#     """
#     if isinstance(data, list):
#         # Duyệt qua từng phần tử trong list để sắp xếp children trước
#         for item in data:
#             sort_json_by_geometry(item)
        
#         # Tiến hành sắp xếp các phần tử cùng cấp trong list
#         # Định nghĩa hàm key để lấy tọa độ an toàn, tránh trường hợp node không có bounds
#         def get_sort_key(node):
#             if isinstance(node, dict) and "bounds" in node:
#                 b = node["bounds"]
#                 # Ưu tiên y trước (làm tròn để tránh designer lệch 1-2px lẻ), sau đó đến x
#                 return (round(b.get("y", 0)), round(b.get("x", 0)))
#             return (0, 0)
        
#         data.sort(key=get_sort_key)

#     elif isinstance(data, dict):
#         # Nếu node có children, đệ quy sâu vào trong để sắp xếp cụm con của nó
#         if "children" in data and isinstance(data["children"], list):
#             sort_json_by_geometry(data["children"])
            
#     return data

def make_it_pretty(raw_result):
    """
    Hàm bóc tách phần text bị escape lồng bên trong response MCP 
    để chuyển đổi ngược thành JSON Object sạch sẽ.
    """
    try:
        # Điều hướng sâu vào cấu trúc: result -> content -> phần tử đầu tiên
        content_list = raw_result.get("result", {}).get("content", [])
        if not content_list:
            return None

        first_content = content_list[0]
        if first_content.get("type") == "text" and "text" in first_content:
            # Bóc tách chuỗi string bị nén ngược lại thành List/Dict trong Python
            clean_data = json.loads(first_content["text"])

            # --- CHUẨN HÓA THỨ TỰ HÌNH HỌC Ở ĐÂY ---
            sorted_data = sort_json_by_geometry(clean_data)
            return sorted_data
    except Exception as e:
        print(f"[Cảnh báo làm đẹp dữ liệu lỗi]: {e}")
    return None


if __name__ == "__main__":
    # --- BƯỚC TEST THỰC TẾ ---
    # Hãy đảm bảo bạn ĐANG SELECT một node/màn hình bất kỳ trên Figma Desktop trước khi chạy

    # Tạo thư mục output nếu chưa tồn tại
    output_dir = "output"
    if not os.path.exists(output_dir):
        os.makedirs(output_dir)

    # 1. Chạy get_selection trước để lấy ID động từ Figma Desktop
    print("=== BƯỚC 1: LẤY NODE ĐANG CHỌN TỪ FIGMA ===")
    selection_result = call_mcp_tool("get_selection", {})

    node_id = None
    if selection_result:
        # Lưu lại file kết quả get_selection gốc
        # with open(os.path.join(output_dir, "get_selection_raw.json"), "w", encoding="utf-8") as f:
        #     json.dump(selection_result, f, indent=4, ensure_ascii=False)
            
        pretty_selection = make_it_pretty(selection_result)
        if pretty_selection and isinstance(pretty_selection, list) and len(pretty_selection) > 0:
            with open(os.path.join(output_dir, "get_selection_pretty.json"), "w", encoding="utf-8") as f:
                json.dump(pretty_selection, f, indent=4, ensure_ascii=False)
                
            node_id = pretty_selection[0].get("id")
            print(f"🎯 Đã tìm thấy Node ID đang chọn: {node_id}\n")
        else:
            print("⚠️ LƯU Ý: Bạn chưa select node nào trên Figma Desktop. Một số tool cần ID sẽ chạy với giá trị mặc định hoặc bỏ qua.")

    # Nếu không lấy được node_id từ selection, thử lấy từ viewport hoặc fallback
    if not node_id:
        print("--> Đang thử lấy thông tin context thay thế từ get_viewport...")
        viewport_res = call_mcp_tool("get_viewport", {})
        # Dự phòng một chuỗi rỗng hợp lệ nếu hoàn toàn không có tương tác
        node_id = "" 

    # Định nghĩa danh sách toàn bộ các tool cần quét và tham số tương ứng của chúng
    # Sử dụng node_id lấy được ở trên cho các tool cần truyền ID
    tools_to_run = [
        # --- Read — Document & Selection ---
        # ("get_document", {}),
        # ("get_metadata", {}),
        # ("get_pages", {}),
        ("get_selection", {}),
        # ("get_viewport", {}),
        # Sửa tham số truyền vào đảm bảo đúng định dạng của figma-mcp-go
        # ("get_node", {"node_id": node_id, "nodeId": node_id} if node_id else {}),
        # ("get_nodes_info", {"node_ids": [node_id], "nodeIds": [node_id]} if node_id else {}),
        # ("get_design_context", {"node_id": node_id, "nodeId": node_id, "detail": "full"} if node_id else {}),
        # Fix lỗi 'query is required': gán mặc định chuỗi tìm kiếm là "a" hoặc ký tự bất kỳ để tránh lỗi trống
        # ("search_nodes", {"node_id": node_id, "nodeId": node_id, "query": "a"} if node_id else {"query": "a"}),
        # ("scan_text_nodes", {"node_id": node_id, "nodeId": node_id} if node_id else {}),
        # ("scan_nodes_by_types", {"node_id": node_id, "nodeId": node_id, "types": ["FRAME", "TEXT"]} if node_id else {}),
        
        # --- Read — Styles & Variables ---
        ("get_styles", {}),
        # ("get_variable_defs", {}),
        # ("get_local_components", {}),
        # ("get_annotations", {}),
        # ("get_fonts", {}),
        # ("get_reactions", {"node_id": node_id, "nodeId": node_id} if node_id else {}),

        # --- Export Tools ---
        ("get_screenshot", {"node_id": node_id, "nodeId": node_id} if node_id else {}),
        # ("save_screenshots", {
        #     "nodes": [node_id], 
        #     "node_ids": [node_id], 
        #     "nodeIds": [node_id], 
        #     "output_dir": output_dir, 
        #     "outputDir": output_dir
        # } if node_id else {}),
        # ("export_frames_to_pdf", {"node_ids": [node_id], "nodeIds": [node_id], "output_path": os.path.join(output_dir, "export_frames.pdf"), "outputPath": os.path.join(output_dir, "export_frames.pdf")} if node_id else {}),        
        # ("export_tokens", {"format": "json"}),
    ]

    print("=== BƯỚC 2: QUÉT TOÀN BỘ CÁC TOOLS ===")
    for tool_name, args in tools_to_run:
        # Nếu tool yêu cầu node_id nhưng bạn chưa chọn gì trên Figma thì bỏ qua tool đó
        if args is None:
            print(f"⏩ Bỏ qua {tool_name} vì chưa lấy được node_id từ Figma Desktop.")
            continue

        result = call_mcp_tool(tool_name, args)

        if result:
            # Check xem kết quả trả về từ server có chứa cờ báo lỗi logic không (isError)
            is_error = result.get("result", {}).get("isError", False)
            
            # # --- TẠM COMMENT không export file raw ---
            # # 1. Xuất file gốc (Dữ liệu thô từ MCP Server) vào thư mục output
            # raw_filename = os.path.join(output_dir, f"{tool_name}_raw.json")
            # with open(raw_filename, "w", encoding="utf-8") as f:
            #     json.dump(result, f, indent=4, ensure_ascii=False)

            if is_error:
                print(f"❌ Tool {tool_name} báo lỗi từ Figma API (Đã lưu log vào _raw.json)")
                continue

            # 2. Xử lý làm đẹp dữ liệu cho con người đọc vào thư mục output
            pretty_data = make_it_pretty(result)

            if pretty_data:
                pretty_filename = os.path.join(output_dir, f"{tool_name}_pretty.json")
                with open(pretty_filename, "w", encoding="utf-8") as f:
                    json.dump(pretty_data, f, indent=4, ensure_ascii=False)
                # --- Tạm comment: print(f"✅ Đã lưu: {tool_name} -> {raw_filename} & {pretty_filename}")
                print(f"✅ Đã lưu: {tool_name} -> {pretty_filename}")
                
            else:
                print(f"⚠️ Đã lưu thô: {tool_name} (Không có text lồng để làm đẹp hoặc rỗng)")
        else:
            print(f"❌ Thất bại khi gọi tool: {tool_name}")

    # Đóng tiến trình server Go một cách an toàn sau khi kết thúc tác vụ
    if _mcp_process and _mcp_process.poll() is None:
        _mcp_process.kill()

    print("\n" + "=" * 50)
    print(f"🎉 HOÀN THÀNH TẤT CẢ! Hãy mở thư mục '{output_dir}/' để xem kết quả.")
    print("=" * 50)