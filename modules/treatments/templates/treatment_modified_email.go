package treatments_templates

// TreatmentModifiedEmailTemplate is an email template used to tell the patient that a treatment was modified by the Psychologist
var TreatmentModifiedEmailTemplate = `<h2>Olá {{ .LikeName }} 😊</h2>
<p>Viemos te informar que o seu tratamento no PSI com {{ .PsyFullName }} foi modificado.</p>
<p>Entre no site para verificar os novos parâmetros do seu tratamento.</p>
<a href="{{ .SiteURL }}">Ir para o site</a>`
