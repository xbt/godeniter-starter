package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

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
	if !strings.Contains(w.Body.String(), "框架特性全景体验中心") {
		t.Errorf("页面未包含特性体验中心标题")
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

	// 5. 创建有效新文章 (POST /admin/articles/create)
	createForm := url.Values{}
	createForm.Set("title", "测试自动创建文章")
	createForm.Set("author", "测试作者")
	createForm.Set("content", "这是自动化单测中发布的内容，字数符合要求。")
	wCreate := httptest.NewRecorder()
	reqCreate, _ := http.NewRequest("POST", "/admin/articles/create", strings.NewReader(createForm.Encode()))
	reqCreate.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqCreate.Header.Set("Cookie", cookie)
	app.ServeHTTP(wCreate, reqCreate)
	if wCreate.Code != http.StatusFound {
		t.Fatalf("创建文章预期重定向到列表，实际: %d", wCreate.Code)
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

	// 7. 编辑新文章 (POST /admin/articles/edit/9)
	editForm := url.Values{}
	editForm.Set("title", "已修改标题测试")
	editForm.Set("author", "测试作者")
	editForm.Set("content", "内容已被编辑更新，保持字数充足。")
	wEdit := httptest.NewRecorder()
	reqEdit, _ := http.NewRequest("POST", "/admin/articles/edit/9", strings.NewReader(editForm.Encode()))
	reqEdit.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqEdit.Header.Set("Cookie", cookie)
	app.ServeHTTP(wEdit, reqEdit)
	if wEdit.Code != http.StatusFound {
		t.Fatalf("编辑文章预期重定向，实际: %d", wEdit.Code)
	}

	// 验证修改生效
	wDetail2 := httptest.NewRecorder()
	reqDetail2, _ := http.NewRequest("GET", "/article/9", nil)
	app.ServeHTTP(wDetail2, reqDetail2)
	if !strings.Contains(wDetail2.Body.String(), "已修改标题测试") {
		t.Errorf("文章详情页未呈现修改后的标题")
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
