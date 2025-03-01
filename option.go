package superhub

// OptionType показывает тип данных, который может принимать динамический параметр.
type OptionType string

const (
	OptionTypeString  OptionType = "STRING"
	OptionTypeInteger OptionType = "INTEGER"
	OptionTypeFloat   OptionType = "FLOAT"
	OptionTypeBoolean OptionType = "BOOLEAN"
)

type OptionPurpose string

const (
	// OptionPurposeNodeIdentifier используется в параметре, отвечающем за идентификатор ноды, на которой должен быть
	// расположен инстанс. Ожидается, что такой параметр будет содержать UUID ноды, который может быть использован,
	// чтобы получить расширенную информацию о ноде.
	OptionPurposeNodeIdentifier OptionPurpose = "NODE_IDENTIFIER"

	// OptionPurposeProcessorLimit используется в параметре, отвечающем за лимит по использованию процессора инстансом.
	// Такой параметр должен иметь числовой тип, а также содержать описание единицы измерения в поле
	// [OptionDescription.Measure]. Ожидается, что значение [Measure.Target] у такого параметра будет равно
	// MeasureTargetProcessorCores.
	OptionPurposeProcessorLimit OptionPurpose = "PROCESSOR_LIMIT"

	// OptionPurposeMemoryLimit используется в параметре, отвечающем за лимит по использованию оперативной памяти
	// инстансом. Такой параметр должен иметь числовой тип, а также содержать описание единицы измерения в поле
	// [OptionDescription.Measure]. Ожидается, что значение [Measure.Target] у такого параметра будет равно
	// MeasureTargetInformationQuantity.
	OptionPurposeMemoryLimit OptionPurpose = "MEMORY_LIMIT"

	// OptionPurposeDiskLimit используется в параметре, отвечающем за лимит по использованию дискового пространства
	// инстансом. Такой параметр должен иметь числовой тип, а также содержать описание единицы измерения в поле
	// [OptionDescription.Measure]. Ожидается, что значение [Measure.Target] у такого параметра будет равно
	// MeasureTargetInformationQuantity.
	OptionPurposeDiskLimit OptionPurpose = "DISK_LIMIT"

	// OptionPurposeOther используется для параметров, для которых не может быть использован ни один из других
	// определённых выше вариантов OptionPurpose.
	OptionPurposeOther OptionPurpose = "OTHER"
)

// Option содержит описание возможного значения для динамического параметра.
type Option struct {
	// Фактическое значение.
	Value string `json:"value"`

	// Отображаемое название параметра.
	DisplayName string `json:"displayName"`

	// Описание параметра.
	Description string `json:"description"`
}

// OptionDescription содержит описание динамического параметра.
type OptionDescription struct {
	// Ключ, который используется для задания параметра.
	Key string `json:"key"`

	// Тип данных, который может принимать данный параметр.
	Type OptionType `json:"type"`

	// Отображаемое название параметра.
	DisplayName string `json:"displayName"`

	// Описание параметра.
	Description string `json:"description"`

	// Предназначение параметра. Возможные значения смотри в OptionPurpose.
	Purpose OptionPurpose `json:"purpose"`

	// Единица измерения, в которой задаётся данный параметр.
	Measure *Measure `json:"measure"`

	// Можно ли изменять параметр после создания сущности, которой этот параметр принадлежит?
	IsEditable bool `json:"isEditable"`

	// Возможные значения параметра.
	Options []Option `json:"options"`
}

// OptionValue содержит значение для динамического параметра, имеющего ключ Key.
type OptionValue struct {
	// Ключ динамического параметра.
	Key string `json:"key"`

	// Значение динамического параметра.
	Value string `json:"value"`
}

// MeasureTarget — измеряемая величина.
type MeasureTarget string

const (
	// MeasureTargetProcessorCores используется для параметров, описывающих количество ядер ЦПУ. Возможные единицы
	// измерения: vCPU, %.
	MeasureTargetProcessorCores MeasureTarget = "PROCESSOR_CORES"

	// MeasureTargetInformationQuantity используется для параметров, описывающих объём информации. Возможные единицы
	// измерения: GB, MB.
	MeasureTargetInformationQuantity MeasureTarget = "INFORMATION_QUANTITY"
)

type Measure struct {
	Target MeasureTarget `json:"target"`
	Unit   string        `json:"unit"`
}
