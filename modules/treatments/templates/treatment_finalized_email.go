package treatments_templates

// TreatmentFinalizedEmailTemplate is an email template used to tell the patient when a treatment is finalized
var TreatmentFinalizedEmailTemplate = `<h2>Olá {{ .LikeName }} 😊</h2>
<p>Viemos te informar que o seu tratamento no PSI com {{ .PsyFullName }} foi finalizado. Isso significa que ela/ele entende que o tratamento foi um sucesso e consultas futuras não são mais necessárias.</p>
<p>Esperamos que a sua experiência usando nossa plataforma tenha sido proveitosa. Seu usuário ainda estará ativo caso deseje iniciar um novo tratamento, ou deixar um elogio, sugestão ou reclamação.</p>
<a href="{{ .SiteURL }}">Ir para o site</a>`
