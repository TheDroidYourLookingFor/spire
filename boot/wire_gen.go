// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package boot

import (
	"github.com/Akkadius/spire/internal/analytics"
	"github.com/Akkadius/spire/internal/app"
	"github.com/Akkadius/spire/internal/assets"
	"github.com/Akkadius/spire/internal/auditlog"
	"github.com/Akkadius/spire/internal/auth"
	"github.com/Akkadius/spire/internal/backup"
	"github.com/Akkadius/spire/internal/clientfiles"
	"github.com/Akkadius/spire/internal/connection"
	"github.com/Akkadius/spire/internal/console/cmd"
	"github.com/Akkadius/spire/internal/database"
	"github.com/Akkadius/spire/internal/deploy"
	"github.com/Akkadius/spire/internal/desktop"
	"github.com/Akkadius/spire/internal/encryption"
	"github.com/Akkadius/spire/internal/eqemuanalytics"
	"github.com/Akkadius/spire/internal/eqemuchangelog"
	"github.com/Akkadius/spire/internal/eqemuserver"
	"github.com/Akkadius/spire/internal/eqemuserverconfig"
	"github.com/Akkadius/spire/internal/eqtraders"
	"github.com/Akkadius/spire/internal/generators"
	"github.com/Akkadius/spire/internal/github"
	"github.com/Akkadius/spire/internal/http"
	"github.com/Akkadius/spire/internal/http/controllers"
	"github.com/Akkadius/spire/internal/http/crudcontrollers"
	"github.com/Akkadius/spire/internal/http/middleware"
	"github.com/Akkadius/spire/internal/http/staticmaps"
	"github.com/Akkadius/spire/internal/influx"
	"github.com/Akkadius/spire/internal/logger"
	"github.com/Akkadius/spire/internal/pathmgmt"
	"github.com/Akkadius/spire/internal/permissions"
	"github.com/Akkadius/spire/internal/query"
	"github.com/Akkadius/spire/internal/questapi"
	"github.com/Akkadius/spire/internal/spire"
	"github.com/Akkadius/spire/internal/system"
	"github.com/Akkadius/spire/internal/telnet"
	"github.com/Akkadius/spire/internal/user"
	"github.com/Akkadius/spire/internal/websocket"
	"github.com/gertd/go-pluralize"
)

import (
	_ "github.com/go-sql-driver/mysql"
)

// Injectors from wire.go:

func InitializeApplication() (App, error) {
	logrusLogger, err := provideLogger()
	if err != nil {
		return App{}, err
	}
	pathManagement := pathmgmt.NewPathManagement(logrusLogger)
	config := eqemuserverconfig.NewConfig(logrusLogger, pathManagement)
	db, err := provideEQEmuLocalDatabase(config)
	if err != nil {
		return App{}, err
	}
	cache := provideCache()
	mysql := backup.NewMysql(logrusLogger, pathManagement)
	helloWorldCommand := cmd.NewHelloWorldCommand(db, logrusLogger, mysql, pathManagement)
	connections := provideAppDbConnections(config, logrusLogger)
	encrypter := encryption.NewEncrypter(logrusLogger, config)
	resolver := database.NewResolver(connections, logrusLogger, encrypter, cache)
	userUser := user.NewUser(resolver, logrusLogger, encrypter, cache)
	createCommand := user.NewCreateCommand(resolver, logrusLogger, encrypter, userUser)
	modelGeneratorCommand := generators.NewModelGeneratorCommand(db, logrusLogger)
	controllerGeneratorCmd := generators.NewControllerGeneratorCommand(db, logrusLogger)
	helloWorldController := controllers.NewHelloWorldController(db, logrusLogger)
	controller := auth.NewController(resolver, logrusLogger, userUser, cache)
	meController := user.NewMeController()
	client := influx.NewClient()
	analyticsController := analytics.NewController(logrusLogger, client, resolver)
	create := connection.NewCreate(resolver, logrusLogger, encrypter)
	check := connection.NewCheck(resolver, logrusLogger, encrypter)
	pluralizeClient := pluralize.NewClient()
	service := permissions.NewService(resolver, cache, logrusLogger, pluralizeClient)
	settings := spire.NewSettings(connections, logrusLogger)
	init := spire.NewInit(connections, config, logrusLogger, settings, cache, encrypter, create, userUser)
	connectionsController := controllers.NewConnectionsController(resolver, logrusLogger, cache, create, check, service, init, userUser)
	sourceDownloader := github.NewGithubSourceDownloader(logrusLogger, cache)
	parseService := questapi.NewParseService(logrusLogger, cache, sourceDownloader)
	examplesGithubSourcer := questapi.NewExamplesGithubSourcer(logrusLogger, cache, sourceDownloader)
	questapiController := questapi.NewController(logrusLogger, parseService, examplesGithubSourcer)
	appController := app.NewController(cache, logrusLogger, init, userUser, settings, resolver)
	queryController := query.NewController(resolver, logrusLogger)
	exporter := clientfiles.NewExporter(logrusLogger)
	importer := clientfiles.NewImporter(logrusLogger)
	clientfilesController := clientfiles.NewController(logrusLogger, exporter, importer, resolver)
	staticMapController := staticmaps.NewStaticMapController(resolver, logrusLogger)
	releases := eqemuanalytics.NewReleases()
	eqemuanalyticsController := eqemuanalytics.NewController(logrusLogger, resolver, releases)
	authedController := eqemuanalytics.NewAuthedController(logrusLogger, resolver)
	changelog := eqemuchangelog.NewChangelog()
	eqemuchangelogController := eqemuchangelog.NewController(logrusLogger, resolver, changelog)
	deployController := deploy.NewDeployController(logrusLogger)
	assetsController := assets.NewController(logrusLogger, resolver)
	permissionsController := permissions.NewController(logrusLogger, resolver, service)
	userController := user.NewController(resolver, logrusLogger, userUser, encrypter)
	settingsController := spire.NewSettingController(resolver, logrusLogger, encrypter, settings)
	telnetClient := telnet.NewClient(logrusLogger)
	eqemuserverClient := eqemuserver.NewClient(telnetClient)
	updater := eqemuserver.NewUpdater(resolver, logrusLogger, config, settings, pathManagement)
	debugLogger := logger.ProvideDebugLogger()
	launcher := eqemuserver.NewLauncher(debugLogger, config, settings, pathManagement, eqemuserverClient)
	eqemuserverController := eqemuserver.NewController(resolver, logrusLogger, eqemuserverClient, config, pathManagement, settings, updater, launcher)
	publicController := eqemuserver.NewPublicController(resolver, logrusLogger, eqemuserverClient, config, pathManagement, settings, updater)
	eqemuserverconfigController := eqemuserverconfig.NewController(logrusLogger, config)
	backupController := backup.NewController(logrusLogger, mysql, pathManagement)
	handler := websocket.NewHandler(logrusLogger, pathManagement)
	websocketController := websocket.NewController(logrusLogger, pathManagement, handler)
	systemController := system.NewController(logrusLogger)
	bootAppControllerGroups := provideControllers(helloWorldController, controller, meController, analyticsController, connectionsController, questapiController, appController, queryController, clientfilesController, staticMapController, eqemuanalyticsController, authedController, eqemuchangelogController, deployController, assetsController, permissionsController, userController, settingsController, eqemuserverController, publicController, eqemuserverconfigController, backupController, websocketController, systemController)
	userEvent := auditlog.NewUserEvent(resolver, logrusLogger, cache)
	aaAbilityController := crudcontrollers.NewAaAbilityController(resolver, logrusLogger, userEvent)
	aaRankController := crudcontrollers.NewAaRankController(resolver, logrusLogger, userEvent)
	aaRankEffectController := crudcontrollers.NewAaRankEffectController(resolver, logrusLogger, userEvent)
	aaRankPrereqController := crudcontrollers.NewAaRankPrereqController(resolver, logrusLogger, userEvent)
	accountController := crudcontrollers.NewAccountController(resolver, logrusLogger, userEvent)
	accountFlagController := crudcontrollers.NewAccountFlagController(resolver, logrusLogger, userEvent)
	accountIpController := crudcontrollers.NewAccountIpController(resolver, logrusLogger, userEvent)
	accountRewardController := crudcontrollers.NewAccountRewardController(resolver, logrusLogger, userEvent)
	adventureDetailController := crudcontrollers.NewAdventureDetailController(resolver, logrusLogger, userEvent)
	adventureMemberController := crudcontrollers.NewAdventureMemberController(resolver, logrusLogger, userEvent)
	adventureStatController := crudcontrollers.NewAdventureStatController(resolver, logrusLogger, userEvent)
	adventureTemplateController := crudcontrollers.NewAdventureTemplateController(resolver, logrusLogger, userEvent)
	adventureTemplateEntryController := crudcontrollers.NewAdventureTemplateEntryController(resolver, logrusLogger, userEvent)
	adventureTemplateEntryFlavorController := crudcontrollers.NewAdventureTemplateEntryFlavorController(resolver, logrusLogger, userEvent)
	alternateCurrencyController := crudcontrollers.NewAlternateCurrencyController(resolver, logrusLogger, userEvent)
	auraController := crudcontrollers.NewAuraController(resolver, logrusLogger, userEvent)
	baseDatumController := crudcontrollers.NewBaseDatumController(resolver, logrusLogger, userEvent)
	blockedSpellController := crudcontrollers.NewBlockedSpellController(resolver, logrusLogger, userEvent)
	bookController := crudcontrollers.NewBookController(resolver, logrusLogger, userEvent)
	botBuffController := crudcontrollers.NewBotBuffController(resolver, logrusLogger, userEvent)
	botCreateCombinationController := crudcontrollers.NewBotCreateCombinationController(resolver, logrusLogger, userEvent)
	botDatumController := crudcontrollers.NewBotDatumController(resolver, logrusLogger, userEvent)
	botGroupController := crudcontrollers.NewBotGroupController(resolver, logrusLogger, userEvent)
	botGroupMemberController := crudcontrollers.NewBotGroupMemberController(resolver, logrusLogger, userEvent)
	botGuildMemberController := crudcontrollers.NewBotGuildMemberController(resolver, logrusLogger, userEvent)
	botHealRotationController := crudcontrollers.NewBotHealRotationController(resolver, logrusLogger, userEvent)
	botHealRotationMemberController := crudcontrollers.NewBotHealRotationMemberController(resolver, logrusLogger, userEvent)
	botHealRotationTargetController := crudcontrollers.NewBotHealRotationTargetController(resolver, logrusLogger, userEvent)
	botInspectMessageController := crudcontrollers.NewBotInspectMessageController(resolver, logrusLogger, userEvent)
	botInventoryController := crudcontrollers.NewBotInventoryController(resolver, logrusLogger, userEvent)
	botOwnerOptionController := crudcontrollers.NewBotOwnerOptionController(resolver, logrusLogger, userEvent)
	botPetBuffController := crudcontrollers.NewBotPetBuffController(resolver, logrusLogger, userEvent)
	botPetController := crudcontrollers.NewBotPetController(resolver, logrusLogger, userEvent)
	botPetInventoryController := crudcontrollers.NewBotPetInventoryController(resolver, logrusLogger, userEvent)
	botSpellCastingChanceController := crudcontrollers.NewBotSpellCastingChanceController(resolver, logrusLogger, userEvent)
	botSpellSettingController := crudcontrollers.NewBotSpellSettingController(resolver, logrusLogger, userEvent)
	botSpellsEntryController := crudcontrollers.NewBotSpellsEntryController(resolver, logrusLogger, userEvent)
	botStanceController := crudcontrollers.NewBotStanceController(resolver, logrusLogger, userEvent)
	botTimerController := crudcontrollers.NewBotTimerController(resolver, logrusLogger, userEvent)
	bugController := crudcontrollers.NewBugController(resolver, logrusLogger, userEvent)
	bugReportController := crudcontrollers.NewBugReportController(resolver, logrusLogger, userEvent)
	buyerController := crudcontrollers.NewBuyerController(resolver, logrusLogger, userEvent)
	charCreateCombinationController := crudcontrollers.NewCharCreateCombinationController(resolver, logrusLogger, userEvent)
	charCreatePointAllocationController := crudcontrollers.NewCharCreatePointAllocationController(resolver, logrusLogger, userEvent)
	charRecipeListController := crudcontrollers.NewCharRecipeListController(resolver, logrusLogger, userEvent)
	characterActivityController := crudcontrollers.NewCharacterActivityController(resolver, logrusLogger, userEvent)
	characterAltCurrencyController := crudcontrollers.NewCharacterAltCurrencyController(resolver, logrusLogger, userEvent)
	characterAlternateAbilityController := crudcontrollers.NewCharacterAlternateAbilityController(resolver, logrusLogger, userEvent)
	characterAuraController := crudcontrollers.NewCharacterAuraController(resolver, logrusLogger, userEvent)
	characterBandolierController := crudcontrollers.NewCharacterBandolierController(resolver, logrusLogger, userEvent)
	characterBindController := crudcontrollers.NewCharacterBindController(resolver, logrusLogger, userEvent)
	characterBuffController := crudcontrollers.NewCharacterBuffController(resolver, logrusLogger, userEvent)
	characterCorpseController := crudcontrollers.NewCharacterCorpseController(resolver, logrusLogger, userEvent)
	characterCorpseItemController := crudcontrollers.NewCharacterCorpseItemController(resolver, logrusLogger, userEvent)
	characterCurrencyController := crudcontrollers.NewCharacterCurrencyController(resolver, logrusLogger, userEvent)
	characterDatumController := crudcontrollers.NewCharacterDatumController(resolver, logrusLogger, userEvent)
	characterDisciplineController := crudcontrollers.NewCharacterDisciplineController(resolver, logrusLogger, userEvent)
	characterEnabledtaskController := crudcontrollers.NewCharacterEnabledtaskController(resolver, logrusLogger, userEvent)
	characterExpModifierController := crudcontrollers.NewCharacterExpModifierController(resolver, logrusLogger, userEvent)
	characterExpeditionLockoutController := crudcontrollers.NewCharacterExpeditionLockoutController(resolver, logrusLogger, userEvent)
	characterInspectMessageController := crudcontrollers.NewCharacterInspectMessageController(resolver, logrusLogger, userEvent)
	characterInstanceSafereturnController := crudcontrollers.NewCharacterInstanceSafereturnController(resolver, logrusLogger, userEvent)
	characterItemRecastController := crudcontrollers.NewCharacterItemRecastController(resolver, logrusLogger, userEvent)
	characterLanguageController := crudcontrollers.NewCharacterLanguageController(resolver, logrusLogger, userEvent)
	characterLeadershipAbilityController := crudcontrollers.NewCharacterLeadershipAbilityController(resolver, logrusLogger, userEvent)
	characterMaterialController := crudcontrollers.NewCharacterMaterialController(resolver, logrusLogger, userEvent)
	characterMemmedSpellController := crudcontrollers.NewCharacterMemmedSpellController(resolver, logrusLogger, userEvent)
	characterPeqzoneFlagController := crudcontrollers.NewCharacterPeqzoneFlagController(resolver, logrusLogger, userEvent)
	characterPetBuffController := crudcontrollers.NewCharacterPetBuffController(resolver, logrusLogger, userEvent)
	characterPetInfoController := crudcontrollers.NewCharacterPetInfoController(resolver, logrusLogger, userEvent)
	characterPetInventoryController := crudcontrollers.NewCharacterPetInventoryController(resolver, logrusLogger, userEvent)
	characterPotionbeltController := crudcontrollers.NewCharacterPotionbeltController(resolver, logrusLogger, userEvent)
	characterSkillController := crudcontrollers.NewCharacterSkillController(resolver, logrusLogger, userEvent)
	characterSpellController := crudcontrollers.NewCharacterSpellController(resolver, logrusLogger, userEvent)
	characterTaskController := crudcontrollers.NewCharacterTaskController(resolver, logrusLogger, userEvent)
	characterTaskTimerController := crudcontrollers.NewCharacterTaskTimerController(resolver, logrusLogger, userEvent)
	characterTributeController := crudcontrollers.NewCharacterTributeController(resolver, logrusLogger, userEvent)
	chatchannelController := crudcontrollers.NewChatchannelController(resolver, logrusLogger, userEvent)
	chatchannelReservedNameController := crudcontrollers.NewChatchannelReservedNameController(resolver, logrusLogger, userEvent)
	commandSubsettingController := crudcontrollers.NewCommandSubsettingController(resolver, logrusLogger, userEvent)
	completedSharedTaskActivityStateController := crudcontrollers.NewCompletedSharedTaskActivityStateController(resolver, logrusLogger, userEvent)
	completedSharedTaskController := crudcontrollers.NewCompletedSharedTaskController(resolver, logrusLogger, userEvent)
	completedSharedTaskMemberController := crudcontrollers.NewCompletedSharedTaskMemberController(resolver, logrusLogger, userEvent)
	completedTaskController := crudcontrollers.NewCompletedTaskController(resolver, logrusLogger, userEvent)
	contentFlagController := crudcontrollers.NewContentFlagController(resolver, logrusLogger, userEvent)
	damageshieldtypeController := crudcontrollers.NewDamageshieldtypeController(resolver, logrusLogger, userEvent)
	dataBucketController := crudcontrollers.NewDataBucketController(resolver, logrusLogger, userEvent)
	dbStrController := crudcontrollers.NewDbStrController(resolver, logrusLogger, userEvent)
	discordWebhookController := crudcontrollers.NewDiscordWebhookController(resolver, logrusLogger, userEvent)
	discoveredItemController := crudcontrollers.NewDiscoveredItemController(resolver, logrusLogger, userEvent)
	doorController := crudcontrollers.NewDoorController(resolver, logrusLogger, userEvent)
	dynamicZoneController := crudcontrollers.NewDynamicZoneController(resolver, logrusLogger, userEvent)
	dynamicZoneMemberController := crudcontrollers.NewDynamicZoneMemberController(resolver, logrusLogger, userEvent)
	dynamicZoneTemplateController := crudcontrollers.NewDynamicZoneTemplateController(resolver, logrusLogger, userEvent)
	eventlogController := crudcontrollers.NewEventlogController(resolver, logrusLogger, userEvent)
	expeditionController := crudcontrollers.NewExpeditionController(resolver, logrusLogger, userEvent)
	expeditionLockoutController := crudcontrollers.NewExpeditionLockoutController(resolver, logrusLogger, userEvent)
	expeditionMemberController := crudcontrollers.NewExpeditionMemberController(resolver, logrusLogger)
	factionAssociationController := crudcontrollers.NewFactionAssociationController(resolver, logrusLogger, userEvent)
	factionBaseDatumController := crudcontrollers.NewFactionBaseDatumController(resolver, logrusLogger, userEvent)
	factionListController := crudcontrollers.NewFactionListController(resolver, logrusLogger, userEvent)
	factionListModController := crudcontrollers.NewFactionListModController(resolver, logrusLogger, userEvent)
	factionValueController := crudcontrollers.NewFactionValueController(resolver, logrusLogger, userEvent)
	fishingController := crudcontrollers.NewFishingController(resolver, logrusLogger, userEvent)
	forageController := crudcontrollers.NewForageController(resolver, logrusLogger, userEvent)
	friendController := crudcontrollers.NewFriendController(resolver, logrusLogger, userEvent)
	globalLootController := crudcontrollers.NewGlobalLootController(resolver, logrusLogger, userEvent)
	gmIpController := crudcontrollers.NewGmIpController(resolver, logrusLogger, userEvent)
	graveyardController := crudcontrollers.NewGraveyardController(resolver, logrusLogger, userEvent)
	gridController := crudcontrollers.NewGridController(resolver, logrusLogger, userEvent)
	gridEntryController := crudcontrollers.NewGridEntryController(resolver, logrusLogger, userEvent)
	groundSpawnController := crudcontrollers.NewGroundSpawnController(resolver, logrusLogger, userEvent)
	groupIdController := crudcontrollers.NewGroupIdController(resolver, logrusLogger, userEvent)
	guildController := crudcontrollers.NewGuildController(resolver, logrusLogger, userEvent)
	guildMemberController := crudcontrollers.NewGuildMemberController(resolver, logrusLogger, userEvent)
	guildRankController := crudcontrollers.NewGuildRankController(resolver, logrusLogger, userEvent)
	guildRelationController := crudcontrollers.NewGuildRelationController(resolver, logrusLogger, userEvent)
	hackerController := crudcontrollers.NewHackerController(resolver, logrusLogger, userEvent)
	horseController := crudcontrollers.NewHorseController(resolver, logrusLogger, userEvent)
	instanceListController := crudcontrollers.NewInstanceListController(resolver, logrusLogger, userEvent)
	instanceListPlayerController := crudcontrollers.NewInstanceListPlayerController(resolver, logrusLogger, userEvent)
	inventoryController := crudcontrollers.NewInventoryController(resolver, logrusLogger, userEvent)
	inventorySnapshotController := crudcontrollers.NewInventorySnapshotController(resolver, logrusLogger, userEvent)
	ipExemptionController := crudcontrollers.NewIpExemptionController(resolver, logrusLogger, userEvent)
	itemController := crudcontrollers.NewItemController(resolver, logrusLogger, userEvent)
	itemTickController := crudcontrollers.NewItemTickController(resolver, logrusLogger, userEvent)
	ldonTrapEntryController := crudcontrollers.NewLdonTrapEntryController(resolver, logrusLogger, userEvent)
	ldonTrapTemplateController := crudcontrollers.NewLdonTrapTemplateController(resolver, logrusLogger, userEvent)
	levelExpModController := crudcontrollers.NewLevelExpModController(resolver, logrusLogger, userEvent)
	lfguildController := crudcontrollers.NewLfguildController(resolver, logrusLogger, userEvent)
	loginAccountController := crudcontrollers.NewLoginAccountController(resolver, logrusLogger, userEvent)
	loginApiTokenController := crudcontrollers.NewLoginApiTokenController(resolver, logrusLogger, userEvent)
	loginServerAdminController := crudcontrollers.NewLoginServerAdminController(resolver, logrusLogger, userEvent)
	loginServerListTypeController := crudcontrollers.NewLoginServerListTypeController(resolver, logrusLogger, userEvent)
	loginWorldServerController := crudcontrollers.NewLoginWorldServerController(resolver, logrusLogger, userEvent)
	logsysCategoryController := crudcontrollers.NewLogsysCategoryController(resolver, logrusLogger, userEvent)
	lootdropController := crudcontrollers.NewLootdropController(resolver, logrusLogger, userEvent)
	lootdropEntryController := crudcontrollers.NewLootdropEntryController(resolver, logrusLogger, userEvent)
	loottableController := crudcontrollers.NewLoottableController(resolver, logrusLogger, userEvent)
	loottableEntryController := crudcontrollers.NewLoottableEntryController(resolver, logrusLogger, userEvent)
	mailController := crudcontrollers.NewMailController(resolver, logrusLogger, userEvent)
	merchantlistController := crudcontrollers.NewMerchantlistController(resolver, logrusLogger, userEvent)
	merchantlistTempController := crudcontrollers.NewMerchantlistTempController(resolver, logrusLogger, userEvent)
	nameFilterController := crudcontrollers.NewNameFilterController(resolver, logrusLogger, userEvent)
	npcEmoteController := crudcontrollers.NewNpcEmoteController(resolver, logrusLogger, userEvent)
	npcFactionController := crudcontrollers.NewNpcFactionController(resolver, logrusLogger, userEvent)
	npcFactionEntryController := crudcontrollers.NewNpcFactionEntryController(resolver, logrusLogger, userEvent)
	npcScaleGlobalBaseController := crudcontrollers.NewNpcScaleGlobalBaseController(resolver, logrusLogger, userEvent)
	npcSpellController := crudcontrollers.NewNpcSpellController(resolver, logrusLogger, userEvent)
	npcSpellsEffectController := crudcontrollers.NewNpcSpellsEffectController(resolver, logrusLogger, userEvent)
	npcSpellsEffectsEntryController := crudcontrollers.NewNpcSpellsEffectsEntryController(resolver, logrusLogger, userEvent)
	npcSpellsEntryController := crudcontrollers.NewNpcSpellsEntryController(resolver, logrusLogger, userEvent)
	npcTypeController := crudcontrollers.NewNpcTypeController(resolver, logrusLogger, userEvent)
	npcTypesTintController := crudcontrollers.NewNpcTypesTintController(resolver, logrusLogger, userEvent)
	objectContentController := crudcontrollers.NewObjectContentController(resolver, logrusLogger, userEvent)
	objectController := crudcontrollers.NewObjectController(resolver, logrusLogger, userEvent)
	perlEventExportSettingController := crudcontrollers.NewPerlEventExportSettingController(resolver, logrusLogger, userEvent)
	petController := crudcontrollers.NewPetController(resolver, logrusLogger, userEvent)
	petitionController := crudcontrollers.NewPetitionController(resolver, logrusLogger, userEvent)
	petsBeastlordDatumController := crudcontrollers.NewPetsBeastlordDatumController(resolver, logrusLogger, userEvent)
	petsEquipmentsetController := crudcontrollers.NewPetsEquipmentsetController(resolver, logrusLogger, userEvent)
	petsEquipmentsetEntryController := crudcontrollers.NewPetsEquipmentsetEntryController(resolver, logrusLogger, userEvent)
	playerEventLogController := crudcontrollers.NewPlayerEventLogController(resolver, logrusLogger, userEvent)
	playerEventLogSettingController := crudcontrollers.NewPlayerEventLogSettingController(resolver, logrusLogger, userEvent)
	playerTitlesetController := crudcontrollers.NewPlayerTitlesetController(resolver, logrusLogger, userEvent)
	questGlobalController := crudcontrollers.NewQuestGlobalController(resolver, logrusLogger, userEvent)
	raidDetailController := crudcontrollers.NewRaidDetailController(resolver, logrusLogger, userEvent)
	raidMemberController := crudcontrollers.NewRaidMemberController(resolver, logrusLogger, userEvent)
	reportController := crudcontrollers.NewReportController(resolver, logrusLogger, userEvent)
	respawnTimeController := crudcontrollers.NewRespawnTimeController(resolver, logrusLogger, userEvent)
	ruleSetController := crudcontrollers.NewRuleSetController(resolver, logrusLogger, userEvent)
	ruleValueController := crudcontrollers.NewRuleValueController(resolver, logrusLogger, userEvent)
	saylinkController := crudcontrollers.NewSaylinkController(resolver, logrusLogger, userEvent)
	serverScheduledEventController := crudcontrollers.NewServerScheduledEventController(resolver, logrusLogger, userEvent)
	sharedTaskActivityStateController := crudcontrollers.NewSharedTaskActivityStateController(resolver, logrusLogger, userEvent)
	sharedTaskController := crudcontrollers.NewSharedTaskController(resolver, logrusLogger, userEvent)
	sharedTaskDynamicZoneController := crudcontrollers.NewSharedTaskDynamicZoneController(resolver, logrusLogger, userEvent)
	sharedTaskMemberController := crudcontrollers.NewSharedTaskMemberController(resolver, logrusLogger, userEvent)
	skillCapController := crudcontrollers.NewSkillCapController(resolver, logrusLogger, userEvent)
	spawn2Controller := crudcontrollers.NewSpawn2Controller(resolver, logrusLogger, userEvent)
	spawnConditionController := crudcontrollers.NewSpawnConditionController(resolver, logrusLogger, userEvent)
	spawnConditionValueController := crudcontrollers.NewSpawnConditionValueController(resolver, logrusLogger, userEvent)
	spawnEventController := crudcontrollers.NewSpawnEventController(resolver, logrusLogger, userEvent)
	spawnentryController := crudcontrollers.NewSpawnentryController(resolver, logrusLogger, userEvent)
	spawngroupController := crudcontrollers.NewSpawngroupController(resolver, logrusLogger, userEvent)
	spellBucketController := crudcontrollers.NewSpellBucketController(resolver, logrusLogger, userEvent)
	spellGlobalController := crudcontrollers.NewSpellGlobalController(resolver, logrusLogger, userEvent)
	spellsNewController := crudcontrollers.NewSpellsNewController(resolver, logrusLogger, userEvent)
	startZoneController := crudcontrollers.NewStartZoneController(resolver, logrusLogger, userEvent)
	startingItemController := crudcontrollers.NewStartingItemController(resolver, logrusLogger, userEvent)
	taskActivityController := crudcontrollers.NewTaskActivityController(resolver, logrusLogger, userEvent)
	taskController := crudcontrollers.NewTaskController(resolver, logrusLogger, userEvent)
	tasksetController := crudcontrollers.NewTasksetController(resolver, logrusLogger, userEvent)
	timerController := crudcontrollers.NewTimerController(resolver, logrusLogger, userEvent)
	titleController := crudcontrollers.NewTitleController(resolver, logrusLogger, userEvent)
	traderController := crudcontrollers.NewTraderController(resolver, logrusLogger, userEvent)
	tradeskillRecipeController := crudcontrollers.NewTradeskillRecipeController(resolver, logrusLogger, userEvent)
	tradeskillRecipeEntryController := crudcontrollers.NewTradeskillRecipeEntryController(resolver, logrusLogger, userEvent)
	trapController := crudcontrollers.NewTrapController(resolver, logrusLogger, userEvent)
	tributeController := crudcontrollers.NewTributeController(resolver, logrusLogger, userEvent)
	tributeLevelController := crudcontrollers.NewTributeLevelController(resolver, logrusLogger, userEvent)
	variableController := crudcontrollers.NewVariableController(resolver, logrusLogger, userEvent)
	veteranRewardTemplateController := crudcontrollers.NewVeteranRewardTemplateController(resolver, logrusLogger, userEvent)
	zoneController := crudcontrollers.NewZoneController(resolver, logrusLogger, userEvent)
	zoneFlagController := crudcontrollers.NewZoneFlagController(resolver, logrusLogger, userEvent)
	zonePointController := crudcontrollers.NewZonePointController(resolver, logrusLogger, userEvent)
	bootCrudControllers := provideCrudControllers(aaAbilityController, aaRankController, aaRankEffectController, aaRankPrereqController, accountController, accountFlagController, accountIpController, accountRewardController, adventureDetailController, adventureMemberController, adventureStatController, adventureTemplateController, adventureTemplateEntryController, adventureTemplateEntryFlavorController, alternateCurrencyController, auraController, baseDatumController, blockedSpellController, bookController, botBuffController, botCreateCombinationController, botDatumController, botGroupController, botGroupMemberController, botGuildMemberController, botHealRotationController, botHealRotationMemberController, botHealRotationTargetController, botInspectMessageController, botInventoryController, botOwnerOptionController, botPetBuffController, botPetController, botPetInventoryController, botSpellCastingChanceController, botSpellSettingController, botSpellsEntryController, botStanceController, botTimerController, bugController, bugReportController, buyerController, charCreateCombinationController, charCreatePointAllocationController, charRecipeListController, characterActivityController, characterAltCurrencyController, characterAlternateAbilityController, characterAuraController, characterBandolierController, characterBindController, characterBuffController, characterCorpseController, characterCorpseItemController, characterCurrencyController, characterDatumController, characterDisciplineController, characterEnabledtaskController, characterExpModifierController, characterExpeditionLockoutController, characterInspectMessageController, characterInstanceSafereturnController, characterItemRecastController, characterLanguageController, characterLeadershipAbilityController, characterMaterialController, characterMemmedSpellController, characterPeqzoneFlagController, characterPetBuffController, characterPetInfoController, characterPetInventoryController, characterPotionbeltController, characterSkillController, characterSpellController, characterTaskController, characterTaskTimerController, characterTributeController, chatchannelController, chatchannelReservedNameController, commandSubsettingController, completedSharedTaskActivityStateController, completedSharedTaskController, completedSharedTaskMemberController, completedTaskController, contentFlagController, damageshieldtypeController, dataBucketController, dbStrController, discordWebhookController, discoveredItemController, doorController, dynamicZoneController, dynamicZoneMemberController, dynamicZoneTemplateController, eventlogController, expeditionController, expeditionLockoutController, expeditionMemberController, factionAssociationController, factionBaseDatumController, factionListController, factionListModController, factionValueController, fishingController, forageController, friendController, globalLootController, gmIpController, graveyardController, gridController, gridEntryController, groundSpawnController, groupIdController, guildController, guildMemberController, guildRankController, guildRelationController, hackerController, horseController, instanceListController, instanceListPlayerController, inventoryController, inventorySnapshotController, ipExemptionController, itemController, itemTickController, ldonTrapEntryController, ldonTrapTemplateController, levelExpModController, lfguildController, loginAccountController, loginApiTokenController, loginServerAdminController, loginServerListTypeController, loginWorldServerController, logsysCategoryController, lootdropController, lootdropEntryController, loottableController, loottableEntryController, mailController, merchantlistController, merchantlistTempController, nameFilterController, npcEmoteController, npcFactionController, npcFactionEntryController, npcScaleGlobalBaseController, npcSpellController, npcSpellsEffectController, npcSpellsEffectsEntryController, npcSpellsEntryController, npcTypeController, npcTypesTintController, objectContentController, objectController, perlEventExportSettingController, petController, petitionController, petsBeastlordDatumController, petsEquipmentsetController, petsEquipmentsetEntryController, playerEventLogController, playerEventLogSettingController, playerTitlesetController, questGlobalController, raidDetailController, raidMemberController, reportController, respawnTimeController, ruleSetController, ruleValueController, saylinkController, serverScheduledEventController, sharedTaskActivityStateController, sharedTaskController, sharedTaskDynamicZoneController, sharedTaskMemberController, skillCapController, spawn2Controller, spawnConditionController, spawnConditionValueController, spawnEventController, spawnentryController, spawngroupController, spellBucketController, spellGlobalController, spellsNewController, startZoneController, startingItemController, taskActivityController, taskController, tasksetController, timerController, titleController, traderController, tradeskillRecipeController, tradeskillRecipeEntryController, trapController, tributeController, tributeLevelController, variableController, veteranRewardTemplateController, zoneController, zoneFlagController, zonePointController)
	contextMiddleware := user.NewContextMiddleware(resolver, cache, logrusLogger)
	readOnlyMiddleware := middleware.NewReadOnlyMiddleware(resolver, logrusLogger)
	permissionsMiddleware := middleware.NewPermissionsMiddleware(resolver, logrusLogger, cache, service)
	requestLogMiddleware := middleware.NewRequestLogMiddleware(client)
	localUserAuthMiddleware := middleware.NewLocalUserAuthMiddleware(resolver, logrusLogger, cache, settings, init)
	spireAssets := assets.NewSpireAssets(logrusLogger, pathManagement)
	router := NewRouter(bootAppControllerGroups, bootCrudControllers, contextMiddleware, readOnlyMiddleware, permissionsMiddleware, requestLogMiddleware, localUserAuthMiddleware, spireAssets)
	questHotReloadWatcher := eqemuserver.NewQuestHotReloadWatcher(debugLogger, config, pathManagement, eqemuserverClient, resolver)
	server := http.NewServer(logrusLogger, router, questHotReloadWatcher)
	httpServeCommand := cmd.NewHttpServeCommand(logrusLogger, server)
	routesListCommand := cmd.NewRoutesListCommand(router, logrusLogger)
	configurationCommand := generators.NewGenerateConfigurationCommand(resolver, logrusLogger)
	migrateCommand := spire.NewMigrateCommand(connections, logrusLogger)
	parseCommand := questapi.NewParseCommand(logrusLogger, parseService)
	exampleTestCommand := questapi.NewExampleTestCommand(logrusLogger, examplesGithubSourcer)
	raceModelMapsCommand := generators.NewRaceModelMapsCommand(logrusLogger)
	changelogCommand := eqemuchangelog.NewChangelogCommand(db, logrusLogger, changelog)
	testFilesystemCommand := cmd.NewTestFilesystemCommand(logrusLogger, pathManagement)
	initCommand := spire.NewInitCommand(logrusLogger, init)
	changePasswordCommand := user.NewChangePasswordCommand(resolver, logrusLogger, encrypter, userUser)
	serverLauncherCommand := spire.NewServerLauncherCommand(logrusLogger, pathManagement)
	crashAnalyticsFingerprintBackfillCommand := spire.NewCrashAnalyticsCommand(logrusLogger, pathManagement, resolver)
	updateCommand := eqemuserver.NewUpdateCommand(logrusLogger, config, settings, pathManagement, launcher, updater)
	launcherCmd := eqemuserver.NewLauncherCmd(logrusLogger, launcher)
	scrapeCommand := eqtraders.NewScrapeCommand(db, logrusLogger)
	importCommand := eqtraders.NewImportCommand(db, logrusLogger)
	v := ProvideCommands(helloWorldCommand, createCommand, modelGeneratorCommand, controllerGeneratorCmd, httpServeCommand, routesListCommand, configurationCommand, migrateCommand, parseCommand, exampleTestCommand, raceModelMapsCommand, changelogCommand, testFilesystemCommand, initCommand, changePasswordCommand, serverLauncherCommand, crashAnalyticsFingerprintBackfillCommand, updateCommand, launcherCmd, scrapeCommand, importCommand)
	webBoot := desktop.NewWebBoot(logrusLogger, server, config)
	bootApp := NewApplication(db, logrusLogger, cache, v, resolver, connections, router, webBoot, init)
	return bootApp, nil
}
