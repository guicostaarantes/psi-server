package templates

// CreateUserEmailTemplate ...
var CreateUserEmailTemplate = `<h2>Oi {{ .FirstName }}</h2>
<p>Você recebeu esse email pois se cadastrou no PSI.</p>
<p>Para confirmarmos que este email é realmente seu, clique no botão abaixo e crie uma senha pra você acessar nosso site.</p>
<a href="{{ .SiteURL }}/cadastro?token={{ .Token }}">Criar uma senha</a>`
