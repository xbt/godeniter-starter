package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"godeniter-starter/app/services"
	"godeniter-starter/config"
)

// TestMiddlewarePipelineAndHeaders 测试全局中间件流水线 (计时与安全防护头注入)
func TestMiddlewarePipelineAndHeaders(t *testing.T) {
	cfg := config.DefaultConfig()
	app := setupApp(cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("首页应返回 200，实际: %d", w.Code)
	}

	// 1. 验证 ResponseTimer 中间件注入的响应头
	if w.Header().Get("X-Response-Time") == "" {
		t.Errorf("预期 ResponseTimer 中间件注入 X-Response-Time 响应头")
	}
	if !strings.Contains(w.Header().Get("Server-Timing"), "app;dur=") {
		t.Errorf("预期 ResponseTimer 中间件注入 Server-Timing 响应头")
	}

	// 2. 验证 SecurityHeaders 中间件注入的安全防护头
	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("预期注入 X-Content-Type-Options: nosniff")
	}
	if w.Header().Get("X-Frame-Options") != "SAMEORIGIN" {
		t.Errorf("预期注入 X-Frame-Options: SAMEORIGIN")
	}
}

// TestFavicon 测试单文件内嵌并正确分发自定义 /favicon.ico 图标
func TestFavicon(t *testing.T) {
	cfg := config.DefaultConfig()
	app := setupApp(cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/favicon.ico", nil)
	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("请求 /favicon.ico 应返回 200，实际: %d", w.Code)
	}
	if w.Header().Get("Content-Type") != "image/x-icon" {
		t.Errorf("Content-Type 预期为 image/x-icon，实际: %s", w.Header().Get("Content-Type"))
	}
	if w.Body.Len() == 0 {
		t.Errorf("favicon.ico 数据长度不应为空")
	}
}


// TestPanicRecovery 测试中间件优雅捕获 Panic 并维持进程不崩服
func TestPanicRecovery(t *testing.T) {
	cfg := config.DefaultConfig()
	app := setupApp(cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/demo/panic", nil)
	app.ServeHTTP(w, req)

	// 预期被 Recovery 中间件拦截并返回 500 内部服务器错误，服务保持可用
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("故意触发 Panic 应被 Recovery 拦截并返回 500，实际返回: %d", w.Code)
	}

	// 再次发起正常请求，确认服务未宕机，依然平稳运行
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/", nil)
	app.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("Panic 触发后服务应继续正常处理后续请求，实际返回: %d", w2.Code)
	}
}

// TestFeaturesPage 测试特性体验中心页面
func TestFeaturesPage(t *testing.T) {
	cfg := config.DefaultConfig()
	app := setupApp(cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/features", nil)
	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("特性中心应返回 200，实际: %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "框架特性全景体验中心") {
		t.Errorf("页面未包含特性体验中心标题")
	}
	// 验证 14 项关键特性卡片均完整呈现
	expectedKeywords := []string{
		"内置秒级与 Linux Crontab 调度引擎",
		"纯 Go 多存储驱动",
		"跨平台桌面托盘与无黑框运行",
		"跨平台守护进程服务化",
		"敏感文件与黑客探测主动防御",
		"强类型反射依赖注入与上下文解耦",
		"纯 Go 0-CGO 嵌入式 SQLite 与链式 ORM",
		"W3C Server-Timing 与安全防护标头",
		"Panic 优雅恢复与故障隔离",
		"无侵入 HTML 注释模板渲染语法",
		"封面上传与静态资源极速分发",
		"零依赖参数校验器",
		"敏感数据脱敏与 XSS 脚本过滤",
		"零依赖 JSON 配置与一键独立打包",
	}
	for _, kw := range expectedKeywords {
		if !strings.Contains(body, kw) {
			t.Errorf("特性页面应包含核心卡片 [%s]", kw)
		}
	}
}

// TestHomePaginationAndSearch 测试前台分页与搜索联动
func TestHomePaginationAndSearch(t *testing.T) {
	cfg := config.DefaultConfig()
	app := setupApp(cfg)

	// 1. 测试第 1 页 (默认 5 条/页，预置 8 篇文章，ID 倒序：8, 7, 6, 5, 4 应在第 1 页)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/", nil)
	app.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Fatalf("预期首页返回 200，实际返回: %d", w1.Code)
	}
	body1 := w1.Body.String()
	if !strings.Contains(body1, "企业级数据安全：敏感信息脱敏与 XSS 过滤") {
		t.Errorf("第 1 页应包含第 8 篇文章标题")
	}
	if !strings.Contains(body1, "ActiveRecord 链式查询构造器实战") {
		t.Errorf("第 1 页应包含第 4 篇文章标题")
	}
	if strings.Contains(body1, "0 依赖轻量依赖注入 (DI) 容器深度剖析") { // 这是第 3 篇文章，应在第 2 页
		t.Errorf("第 1 页不应包含第 3 篇文章标题（单页限制 5 条）")
	}

	// 2. 测试第 2 页 (应包含 ID 3, 2, 1)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/?page=2", nil)
	app.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("预期第 2 页返回 200，实际返回: %d", w2.Code)
	}
	body2 := w2.Body.String()
	if !strings.Contains(body2, "0 依赖轻量依赖注入 (DI) 容器深度剖析") {
		t.Errorf("第 2 页应包含第 3 篇文章标题")
	}
	if !strings.Contains(body2, "欢迎使用 Godeniter 2.0 框架") {
		t.Errorf("第 2 页应包含第 1 篇文章标题")
	}

	// 3. 测试搜索关键词联动
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/?keyword=微服务", nil)
	app.ServeHTTP(w3, req3)
	if w3.Code != http.StatusOK {
		t.Fatalf("预期搜索请求返回 200，实际返回: %d", w3.Code)
	}
	body3 := w3.Body.String()
	if !strings.Contains(body3, "微服务环境下的动态配置与 Sidecar 机制") {
		t.Errorf("搜索 '微服务' 应包含相关文章")
	}
	if strings.Contains(body3, "服务端 Session 会话管理与闪存消息") {
		t.Errorf("搜索 '微服务' 不应包含无关文章")
	}
}

// TestArticleDetail 测试文章详情与 404
func TestArticleDetail(t *testing.T) {
	cfg := config.DefaultConfig()
	app := setupApp(cfg)

	// 1. 访问正常详情页
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/article/1", nil)
	app.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Fatalf("详情页应返回 200，实际: %d", w1.Code)
	}
	if !strings.Contains(w1.Body.String(), "欢迎使用 Godeniter 2.0 框架") {
		t.Errorf("详情页应包含标题")
	}

	// 2. 测试已登录状态下，无侵入注释语法 <!--{{ if .CurrentUser }}--> 正确渲染编辑按钮
	wLogin := httptest.NewRecorder()
	loginData := url.Values{"username": {"admin"}, "password": {"123456"}}
	reqLogin, _ := http.NewRequest("POST", "/login", strings.NewReader(loginData.Encode()))
	reqLogin.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	app.ServeHTTP(wLogin, reqLogin)
	cookie := wLogin.Header().Get("Set-Cookie")

	wAuth := httptest.NewRecorder()
	reqAuth, _ := http.NewRequest("GET", "/article/1", nil)
	reqAuth.Header.Set("Cookie", cookie)
	app.ServeHTTP(wAuth, reqAuth)
	if !strings.Contains(wAuth.Body.String(), `/admin/articles/edit/1`) {
		t.Errorf("注释语法 <!--{{ .Article.ID }}--> 预期在登录后渲染编辑按钮，实际未找到")
	}

	// 3. 访问不存在的详情页
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/article/9999", nil)
	app.ServeHTTP(w2, req2)
	if w2.Code != http.StatusNotFound {
		t.Fatalf("不存在的文章应返回 404，实际: %d", w2.Code)
	}
}


// TestFileUploadAPI 测试 RESTful 文件上传接口与安全限制
func TestFileUploadAPI(t *testing.T) {
	cfg := config.DefaultConfig()
	app := setupApp(cfg)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test_avatar.png")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("fake image binary content for test"))
	writer.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("文件上传应返回 200，实际: %d, body: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "saved_path") {
		t.Errorf("上传响应应包含 saved_path 字段")
	}
}

// TestAdminAuthAndCRUD 测试后台中间件路由守卫、表单校验拦截与完整 CRUD
func TestAdminAuthAndCRUD(t *testing.T) {
	cfg := config.DefaultConfig()
	app := setupApp(cfg)

	// 1. 未登录访问后台，预期被 AuthRequired 中间件拦截并重定向至 /login
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/admin/articles", nil)
	app.ServeHTTP(w1, req1)
	if w1.Code != http.StatusFound {
		t.Fatalf("未登录访问后台预期 302 重定向，实际: %d", w1.Code)
	}
	if w1.Header().Get("Location") != "/login" {
		t.Fatalf("预期重定向至 /login，实际 Location: %s", w1.Header().Get("Location"))
	}

	// 2. 执行登录 (POST /login)
	formData := url.Values{}
	formData.Set("username", "admin")
	formData.Set("password", "123456")
	wLogin := httptest.NewRecorder()
	reqLogin, _ := http.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))
	reqLogin.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	app.ServeHTTP(wLogin, reqLogin)

	if wLogin.Code != http.StatusFound {
		t.Fatalf("登录成功后预期 302 重定向，实际: %d", wLogin.Code)
	}
	cookie := wLogin.Header().Get("Set-Cookie")
	if cookie == "" {
		t.Fatalf("登录成功后预期返回 Set-Cookie")
	}

	// 3. 携带 Cookie 访问后台管理首页
	wAdmin := httptest.NewRecorder()
	reqAdmin, _ := http.NewRequest("GET", "/admin/articles", nil)
	reqAdmin.Header.Set("Cookie", cookie)
	app.ServeHTTP(wAdmin, reqAdmin)
	if wAdmin.Code != http.StatusOK {
		t.Fatalf("携带有效 Session 预期返回 200，实际: %d", wAdmin.Code)
	}
	if !strings.Contains(wAdmin.Body.String(), "文章管理中心") {
		t.Errorf("后台页面应包含'文章管理中心'")
	}

	// 4. 测试表单 Tag 校验拦截 (提交标题只有 1 个字符，未达 min=3 规则)
	badForm := url.Values{}
	badForm.Set("title", "a") // 太短
	badForm.Set("author", "测试作者")
	badForm.Set("content", "内容足够长但是标题太短了")
	wBad := httptest.NewRecorder()
	reqBad, _ := http.NewRequest("POST", "/admin/articles/create", strings.NewReader(badForm.Encode()))
	reqBad.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqBad.Header.Set("Cookie", cookie)
	app.ServeHTTP(wBad, reqBad)
	if !strings.Contains(wBad.Body.String(), "表单参数校验失败") {
		t.Errorf("预期触发表单校验失败提示，实际未拦截")
	}

	// 5. 创建有效新文章 (POST /admin/articles/create，采用页面真实使用的 multipart/form-data 并附带封面上传)
	createBody := &bytes.Buffer{}
	createWriter := multipart.NewWriter(createBody)
	_ = createWriter.WriteField("title", "测试自动创建文章")
	_ = createWriter.WriteField("author", "测试作者")
	_ = createWriter.WriteField("content", "这是自动化单测中发布的内容，字数符合要求。")
	partCreate, _ := createWriter.CreateFormFile("cover", "init_cover.png")
	_, _ = partCreate.Write([]byte("fake png binary data for initial cover"))
	_ = createWriter.Close()

	wCreate := httptest.NewRecorder()
	reqCreate, _ := http.NewRequest("POST", "/admin/articles/create", createBody)
	reqCreate.Header.Set("Content-Type", createWriter.FormDataContentType())
	reqCreate.Header.Set("Cookie", cookie)
	app.ServeHTTP(wCreate, reqCreate)
	if wCreate.Code != http.StatusFound {
		t.Fatalf("创建文章预期重定向到列表，实际: %d, body: %s", wCreate.Code, wCreate.Body.String())
	}

	// 6. 验证新创建的文章在详情页可以访问 (预置 8 篇，新文章 ID 应为 9)
	wDetail := httptest.NewRecorder()
	reqDetail, _ := http.NewRequest("GET", "/article/9", nil)
	app.ServeHTTP(wDetail, reqDetail)
	if wDetail.Code != http.StatusOK {
		t.Fatalf("新文章 #9 详情页应返回 200，实际: %d", wDetail.Code)
	}
	if !strings.Contains(wDetail.Body.String(), "测试自动创建文章") {
		t.Errorf("详情页未找到新创建的文章标题")
	}

	// 7. 编辑新文章 (POST /admin/articles/edit/9，使用 multipart/form-data 模拟带图片上传)
	editBody := &bytes.Buffer{}
	mpWriter := multipart.NewWriter(editBody)
	_ = mpWriter.WriteField("title", "已修改标题测试(带封面上传)")
	_ = mpWriter.WriteField("author", "测试作者")
	_ = mpWriter.WriteField("content", "内容已被编辑更新，并且成功上传了新封面！")
	partCover, _ := mpWriter.CreateFormFile("cover", "new_cover.png")
	_, _ = partCover.Write([]byte("fake png binary data for cover update"))
	_ = mpWriter.Close()

	wEdit := httptest.NewRecorder()
	reqEdit, _ := http.NewRequest("POST", "/admin/articles/edit/9", editBody)
	reqEdit.Header.Set("Content-Type", mpWriter.FormDataContentType())
	reqEdit.Header.Set("Cookie", cookie)
	app.ServeHTTP(wEdit, reqEdit)
	if wEdit.Code != http.StatusFound {
		t.Fatalf("编辑文章预期 302 重定向，实际: %d, body: %s", wEdit.Code, wEdit.Body.String())
	}

	// 验证修改生效且封面已被更新
	wDetail2 := httptest.NewRecorder()
	reqDetail2, _ := http.NewRequest("GET", "/article/9", nil)
	app.ServeHTTP(wDetail2, reqDetail2)
	if !strings.Contains(wDetail2.Body.String(), "已修改标题测试(带封面上传)") {
		t.Errorf("文章详情页未呈现修改后的标题")
	}
	if !strings.Contains(wDetail2.Body.String(), "/uploads/images/") {
		t.Errorf("文章详情页未呈现新上传的封面图片路径")
	}

	// 8. 删除新文章 (GET /admin/articles/delete/9)
	wDel := httptest.NewRecorder()
	reqDel, _ := http.NewRequest("GET", "/admin/articles/delete/9", nil)
	reqDel.Header.Set("Cookie", cookie)
	app.ServeHTTP(wDel, reqDel)
	if wDel.Code != http.StatusFound {
		t.Fatalf("删除文章预期重定向，实际: %d", wDel.Code)
	}

	// 验证删除后 404
	wDetail3 := httptest.NewRecorder()
	reqDetail3, _ := http.NewRequest("GET", "/article/9", nil)
	app.ServeHTTP(wDetail3, reqDetail3)
	if wDetail3.Code != http.StatusNotFound {
		t.Fatalf("删除后访问文章应返回 404，实际: %d", wDetail3.Code)
	}
}

// TestBlockSensitiveFilesProbe 验证敏感系统文件与数据库探测拦截中间件 (403 阻断)
func TestBlockSensitiveFilesProbe(t *testing.T) {
	cfg := config.DefaultConfig()
	app := setupApp(cfg)

	// 探测 SQLite 库或 config.json 预期直接 403 Forbidden
	probes := []string{"/data/app.db", "/config.json", "/.env", "/.git/config"}
	for _, p := range probes {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", p, nil)
		app.ServeHTTP(w, req)
		if w.Code != http.StatusForbidden {
			t.Errorf("探测敏感路径 %s 预期返回 403 Forbidden，实际: %d", p, w.Code)
		}
	}
}

// TestStarterStorageDriver 验证演示的多存储驱动接口与本地驱动基本功能
func TestStarterStorageDriver(t *testing.T) {
	tmpDir := t.TempDir()
	driver := services.NewDefaultStorageDriver(tmpDir)

	if driver.Name() != "local" {
		t.Fatalf("预期默认存储驱动为 'local', 实际: %s", driver.Name())
	}

	testData := []byte("fake starter image data")
	fileURL, err := driver.Save("test_sample.png", bytes.NewReader(testData), int64(len(testData)), "image/png")
	if err != nil {
		t.Fatalf("保存文件失败: %v", err)
	}

	if !strings.HasPrefix(fileURL, "/uploads/images/") {
		t.Fatalf("返回 URL 格式异常: %s", fileURL)
	}

	if err := driver.Delete("test_sample.png"); err != nil {
		t.Fatalf("删除文件失败: %v", err)
	}
}

// TestStarterCronScheduler 验证 starter 中内置 Cron 计划任务的注册与状态查询接口
func TestStarterCronScheduler(t *testing.T) {
	cfg := config.DefaultConfig()
	app := setupApp(cfg)

	// 1. 验证 Cron 实例中已注册 heartbeat 与 cleanup
	jobs := app.Cron.Jobs()
	if len(jobs) < 2 {
		t.Fatalf("预期至少注册 2 个计划任务，实际: %d", len(jobs))
	}

	foundHeartbeat := false
	for _, j := range jobs {
		if j.ID == "heartbeat" {
			foundHeartbeat = true
			if j.Spec != "*/15 * * * * *" {
				t.Errorf("heartbeat 任务规格不匹配: %s", j.Spec)
			}
		}
	}
	if !foundHeartbeat {
		t.Errorf("未找到预期的 heartbeat 任务")
	}

	// 2. 通过 REST API 接口获取任务快照
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cron/jobs", nil)
	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("访问 /api/v1/cron/jobs 应返回 200，实际: %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "heartbeat") {
		t.Errorf("API 未返回 heartbeat 任务")
	}

	// 3. 测试通过 API 手动立即触发任务
	wTrigger := httptest.NewRecorder()
	reqTrigger, _ := http.NewRequest("POST", "/api/v1/cron/trigger/heartbeat", nil)
	app.ServeHTTP(wTrigger, reqTrigger)
	if wTrigger.Code != http.StatusOK {
		t.Fatalf("预期触发任务返回 200，实际: %d", wTrigger.Code)
	}

	// 4. 测试触发不存在的任务返回业务错误码 400
	wErr := httptest.NewRecorder()
	reqErr, _ := http.NewRequest("POST", "/api/v1/cron/trigger/not_exist_job", nil)
	app.ServeHTTP(wErr, reqErr)
	if !strings.Contains(wErr.Body.String(), `"code":400`) {
		t.Fatalf("触发不存在的任务预期返回业务错误码 400，实际响应: %s", wErr.Body.String())
	}
}

// TestStarterStorageStatusAPI 验证存储驱动状态端点
func TestStarterStorageStatusAPI(t *testing.T) {
	cfg := config.DefaultConfig()
	app := setupApp(cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/storage/status", nil)
	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("预期返回 200，实际: %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "current_driver") || !strings.Contains(body, "supported_drivers") {
		t.Errorf("响应缺少存储驱动字段: %s", body)
	}
}


