package templates

// ResetPasswordEmailTemplate ...
var ResetPasswordEmailTemplate = `<h2>Oi {{ .FirstName }}</h2>
<p>Você recebeu esse email pois pediu uma mudança de senha no PSI.</p>
<p>Para confirmarmos que este email é realmente seu, clique no botão abaixo e crie uma nova senha pra você acessar nosso site.</p>
<a href="{{ .SiteURL }}/nova-senha?token={{ .Token }}">Mudar minha senha</a>`
