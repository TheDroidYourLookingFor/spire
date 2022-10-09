// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package boot

import (
	"github.com/Akkadius/spire/internal/clientfiles"
	"github.com/Akkadius/spire/internal/connection"
	"github.com/Akkadius/spire/internal/console/cmd"
	"github.com/Akkadius/spire/internal/database"
	"github.com/Akkadius/spire/internal/desktop"
	"github.com/Akkadius/spire/internal/encryption"
	"github.com/Akkadius/spire/internal/github"
	"github.com/Akkadius/spire/internal/http/controllers"
	"github.com/Akkadius/spire/internal/http/crudcontrollers"
	"github.com/Akkadius/spire/internal/http/middleware"
	"github.com/Akkadius/spire/internal/http/staticmaps"
	"github.com/Akkadius/spire/internal/influx"
	"github.com/Akkadius/spire/internal/pathmgmt"
	"github.com/Akkadius/spire/internal/permissions"
	"github.com/Akkadius/spire/internal/questapi"
	"github.com/Akkadius/spire/internal/serverconfig"
)

import (
	_ "github.com/go-sql-driver/mysql"
)

// Injectors from wire.go:

func InitializeApplication() (App, error) {
	logger, err := provideLogger()
	if err != nil {
		return App{}, err
	}
	pathManagement := pathmgmt.NewPathManagement(logger)
	eqEmuServerConfig := serverconfig.NewEQEmuServerConfig(logger, pathManagement)
	db, err := provideEQEmuLocalDatabase(eqEmuServerConfig)
	if err != nil {
		return App{}, err
	}
	cache := provideCache()
	helloWorldCommand := cmd.NewHelloWorldCommand(db, logger)
	generateModelsCommand := cmd.NewGenerateModelsCommand(db, logger)
	generateControllersCommand := cmd.NewGenerateControllersCommand(db, logger)
	helloWorldController := controllers.NewHelloWorldController(db, logger)
	connections := provideAppDbConnections(eqEmuServerConfig)
	encrypter := encryption.NewEncrypter()
	databaseResolver := database.NewDatabaseResolver(connections, logger, encrypter, cache)
	authController := controllers.NewAuthController(databaseResolver, logger)
	meController := controllers.NewMeController()
	client := influx.NewClient()
	analyticsController := controllers.NewAnalyticsController(logger, client, databaseResolver)
	dbConnectionCreateService := connection.NewDbConnectionCreateService(databaseResolver, logger, encrypter)
	dbConnectionCheckService := connection.NewDbConnectionCheckService(databaseResolver, logger, encrypter)
	connectionsController := controllers.NewConnectionsController(databaseResolver, logger, cache, dbConnectionCreateService, dbConnectionCheckService)
	docsController := controllers.NewDocsController(databaseResolver, logger)
	githubSourceDownloader := github.NewGithubSourceDownloader(logger, cache)
	parseService := questapi.NewParseService(logger, cache, githubSourceDownloader)
	questExamplesGithubSourcer := questapi.NewQuestExamplesGithubSourcer(logger, cache, githubSourceDownloader)
	questApiController := controllers.NewQuestApiController(logger, parseService, questExamplesGithubSourcer)
	appController := controllers.NewAppController(cache, logger)
	queryController := controllers.NewQueryController(databaseResolver, logger)
	questFileApiController := controllers.NewQuestFileApiController(logger)
	exporter := clientfiles.NewExporter(logger)
	importer := clientfiles.NewImporter(logger)
	clientFilesController := controllers.NewClientFilesController(logger, exporter, importer, databaseResolver)
	staticMapController := staticmaps.NewStaticMapController(databaseResolver, logger)
	assetsController := controllers.NewAssetsController(logger, databaseResolver)
	resourceList := permissions.NewResourceList(logger)
	permissionsController := controllers.NewPermissionsController(logger, databaseResolver, resourceList)
	usersController := controllers.NewUsersController(databaseResolver, logger)
	bootAppControllerGroups := provideControllers(helloWorldController, authController, meController, analyticsController, connectionsController, docsController, questApiController, appController, queryController, questFileApiController, clientFilesController, staticMapController, assetsController, permissionsController, usersController)
	aaAbilityController := crudcontrollers.NewAaAbilityController(databaseResolver, logger)
	aaRankController := crudcontrollers.NewAaRankController(databaseResolver, logger)
	aaRankEffectController := crudcontrollers.NewAaRankEffectController(databaseResolver, logger)
	aaRankPrereqController := crudcontrollers.NewAaRankPrereqController(databaseResolver, logger)
	accountController := crudcontrollers.NewAccountController(databaseResolver, logger)
	accountFlagController := crudcontrollers.NewAccountFlagController(databaseResolver, logger)
	accountIpController := crudcontrollers.NewAccountIpController(databaseResolver, logger)
	accountRewardController := crudcontrollers.NewAccountRewardController(databaseResolver, logger)
	adventureDetailController := crudcontrollers.NewAdventureDetailController(databaseResolver, logger)
	adventureMemberController := crudcontrollers.NewAdventureMemberController(databaseResolver, logger)
	adventureStatController := crudcontrollers.NewAdventureStatController(databaseResolver, logger)
	adventureTemplateController := crudcontrollers.NewAdventureTemplateController(databaseResolver, logger)
	adventureTemplateEntryController := crudcontrollers.NewAdventureTemplateEntryController(databaseResolver, logger)
	adventureTemplateEntryFlavorController := crudcontrollers.NewAdventureTemplateEntryFlavorController(databaseResolver, logger)
	alternateCurrencyController := crudcontrollers.NewAlternateCurrencyController(databaseResolver, logger)
	auraController := crudcontrollers.NewAuraController(databaseResolver, logger)
	baseDatumController := crudcontrollers.NewBaseDatumController(databaseResolver, logger)
	blockedSpellController := crudcontrollers.NewBlockedSpellController(databaseResolver, logger)
	bookController := crudcontrollers.NewBookController(databaseResolver, logger)
	botDatumController := crudcontrollers.NewBotDatumController(databaseResolver, logger)
	bugController := crudcontrollers.NewBugController(databaseResolver, logger)
	bugReportController := crudcontrollers.NewBugReportController(databaseResolver, logger)
	buyerController := crudcontrollers.NewBuyerController(databaseResolver, logger)
	charCreateCombinationController := crudcontrollers.NewCharCreateCombinationController(databaseResolver, logger)
	charCreatePointAllocationController := crudcontrollers.NewCharCreatePointAllocationController(databaseResolver, logger)
	charRecipeListController := crudcontrollers.NewCharRecipeListController(databaseResolver, logger)
	characterActivityController := crudcontrollers.NewCharacterActivityController(databaseResolver, logger)
	characterAltCurrencyController := crudcontrollers.NewCharacterAltCurrencyController(databaseResolver, logger)
	characterAlternateAbilityController := crudcontrollers.NewCharacterAlternateAbilityController(databaseResolver, logger)
	characterAuraController := crudcontrollers.NewCharacterAuraController(databaseResolver, logger)
	characterBandolierController := crudcontrollers.NewCharacterBandolierController(databaseResolver, logger)
	characterBindController := crudcontrollers.NewCharacterBindController(databaseResolver, logger)
	characterBuffController := crudcontrollers.NewCharacterBuffController(databaseResolver, logger)
	characterCorpseController := crudcontrollers.NewCharacterCorpseController(databaseResolver, logger)
	characterCorpseItemController := crudcontrollers.NewCharacterCorpseItemController(databaseResolver, logger)
	characterCurrencyController := crudcontrollers.NewCharacterCurrencyController(databaseResolver, logger)
	characterDatumController := crudcontrollers.NewCharacterDatumController(databaseResolver, logger)
	characterDisciplineController := crudcontrollers.NewCharacterDisciplineController(databaseResolver, logger)
	characterEnabledtaskController := crudcontrollers.NewCharacterEnabledtaskController(databaseResolver, logger)
	characterExpModifierController := crudcontrollers.NewCharacterExpModifierController(databaseResolver, logger)
	characterExpeditionLockoutController := crudcontrollers.NewCharacterExpeditionLockoutController(databaseResolver, logger)
	characterInspectMessageController := crudcontrollers.NewCharacterInspectMessageController(databaseResolver, logger)
	characterInstanceSafereturnController := crudcontrollers.NewCharacterInstanceSafereturnController(databaseResolver, logger)
	characterItemRecastController := crudcontrollers.NewCharacterItemRecastController(databaseResolver, logger)
	characterLanguageController := crudcontrollers.NewCharacterLanguageController(databaseResolver, logger)
	characterLeadershipAbilityController := crudcontrollers.NewCharacterLeadershipAbilityController(databaseResolver, logger)
	characterMaterialController := crudcontrollers.NewCharacterMaterialController(databaseResolver, logger)
	characterMemmedSpellController := crudcontrollers.NewCharacterMemmedSpellController(databaseResolver, logger)
	characterPeqzoneFlagController := crudcontrollers.NewCharacterPeqzoneFlagController(databaseResolver, logger)
	characterPetBuffController := crudcontrollers.NewCharacterPetBuffController(databaseResolver, logger)
	characterPetInfoController := crudcontrollers.NewCharacterPetInfoController(databaseResolver, logger)
	characterPetInventoryController := crudcontrollers.NewCharacterPetInventoryController(databaseResolver, logger)
	characterPotionbeltController := crudcontrollers.NewCharacterPotionbeltController(databaseResolver, logger)
	characterSkillController := crudcontrollers.NewCharacterSkillController(databaseResolver, logger)
	characterSpellController := crudcontrollers.NewCharacterSpellController(databaseResolver, logger)
	characterTaskController := crudcontrollers.NewCharacterTaskController(databaseResolver, logger)
	characterTaskTimerController := crudcontrollers.NewCharacterTaskTimerController(databaseResolver, logger)
	completedSharedTaskActivityStateController := crudcontrollers.NewCompletedSharedTaskActivityStateController(databaseResolver, logger)
	completedSharedTaskController := crudcontrollers.NewCompletedSharedTaskController(databaseResolver, logger)
	completedSharedTaskMemberController := crudcontrollers.NewCompletedSharedTaskMemberController(databaseResolver, logger)
	completedTaskController := crudcontrollers.NewCompletedTaskController(databaseResolver, logger)
	contentFlagController := crudcontrollers.NewContentFlagController(databaseResolver, logger)
	damageshieldtypeController := crudcontrollers.NewDamageshieldtypeController(databaseResolver, logger)
	dataBucketController := crudcontrollers.NewDataBucketController(databaseResolver, logger)
	dbStrController := crudcontrollers.NewDbStrController(databaseResolver, logger)
	discordWebhookController := crudcontrollers.NewDiscordWebhookController(databaseResolver, logger)
	discoveredItemController := crudcontrollers.NewDiscoveredItemController(databaseResolver, logger)
	doorController := crudcontrollers.NewDoorController(databaseResolver, logger)
	dynamicZoneController := crudcontrollers.NewDynamicZoneController(databaseResolver, logger)
	dynamicZoneMemberController := crudcontrollers.NewDynamicZoneMemberController(databaseResolver, logger)
	dynamicZoneTemplateController := crudcontrollers.NewDynamicZoneTemplateController(databaseResolver, logger)
	eventlogController := crudcontrollers.NewEventlogController(databaseResolver, logger)
	expeditionController := crudcontrollers.NewExpeditionController(databaseResolver, logger)
	expeditionLockoutController := crudcontrollers.NewExpeditionLockoutController(databaseResolver, logger)
	expeditionMemberController := crudcontrollers.NewExpeditionMemberController(databaseResolver, logger)
	factionAssociationController := crudcontrollers.NewFactionAssociationController(databaseResolver, logger)
	factionBaseDatumController := crudcontrollers.NewFactionBaseDatumController(databaseResolver, logger)
	factionListController := crudcontrollers.NewFactionListController(databaseResolver, logger)
	factionListModController := crudcontrollers.NewFactionListModController(databaseResolver, logger)
	factionValueController := crudcontrollers.NewFactionValueController(databaseResolver, logger)
	fishingController := crudcontrollers.NewFishingController(databaseResolver, logger)
	forageController := crudcontrollers.NewForageController(databaseResolver, logger)
	friendController := crudcontrollers.NewFriendController(databaseResolver, logger)
	globalLootController := crudcontrollers.NewGlobalLootController(databaseResolver, logger)
	gmIpController := crudcontrollers.NewGmIpController(databaseResolver, logger)
	graveyardController := crudcontrollers.NewGraveyardController(databaseResolver, logger)
	gridController := crudcontrollers.NewGridController(databaseResolver, logger)
	gridEntryController := crudcontrollers.NewGridEntryController(databaseResolver, logger)
	groundSpawnController := crudcontrollers.NewGroundSpawnController(databaseResolver, logger)
	groupIdController := crudcontrollers.NewGroupIdController(databaseResolver, logger)
	guildController := crudcontrollers.NewGuildController(databaseResolver, logger)
	guildMemberController := crudcontrollers.NewGuildMemberController(databaseResolver, logger)
	guildRankController := crudcontrollers.NewGuildRankController(databaseResolver, logger)
	guildRelationController := crudcontrollers.NewGuildRelationController(databaseResolver, logger)
	hackerController := crudcontrollers.NewHackerController(databaseResolver, logger)
	horseController := crudcontrollers.NewHorseController(databaseResolver, logger)
	instanceListController := crudcontrollers.NewInstanceListController(databaseResolver, logger)
	instanceListPlayerController := crudcontrollers.NewInstanceListPlayerController(databaseResolver, logger)
	inventoryController := crudcontrollers.NewInventoryController(databaseResolver, logger)
	inventorySnapshotController := crudcontrollers.NewInventorySnapshotController(databaseResolver, logger)
	ipExemptionController := crudcontrollers.NewIpExemptionController(databaseResolver, logger)
	itemController := crudcontrollers.NewItemController(databaseResolver, logger)
	itemTickController := crudcontrollers.NewItemTickController(databaseResolver, logger)
	ldonTrapEntryController := crudcontrollers.NewLdonTrapEntryController(databaseResolver, logger)
	ldonTrapTemplateController := crudcontrollers.NewLdonTrapTemplateController(databaseResolver, logger)
	levelExpModController := crudcontrollers.NewLevelExpModController(databaseResolver, logger)
	lfguildController := crudcontrollers.NewLfguildController(databaseResolver, logger)
	loginAccountController := crudcontrollers.NewLoginAccountController(databaseResolver, logger)
	loginApiTokenController := crudcontrollers.NewLoginApiTokenController(databaseResolver, logger)
	loginServerAdminController := crudcontrollers.NewLoginServerAdminController(databaseResolver, logger)
	loginServerListTypeController := crudcontrollers.NewLoginServerListTypeController(databaseResolver, logger)
	loginWorldServerController := crudcontrollers.NewLoginWorldServerController(databaseResolver, logger)
	logsysCategoryController := crudcontrollers.NewLogsysCategoryController(databaseResolver, logger)
	lootdropController := crudcontrollers.NewLootdropController(databaseResolver, logger)
	lootdropEntryController := crudcontrollers.NewLootdropEntryController(databaseResolver, logger)
	loottableController := crudcontrollers.NewLoottableController(databaseResolver, logger)
	loottableEntryController := crudcontrollers.NewLoottableEntryController(databaseResolver, logger)
	mailController := crudcontrollers.NewMailController(databaseResolver, logger)
	merchantlistController := crudcontrollers.NewMerchantlistController(databaseResolver, logger)
	merchantlistTempController := crudcontrollers.NewMerchantlistTempController(databaseResolver, logger)
	nameFilterController := crudcontrollers.NewNameFilterController(databaseResolver, logger)
	npcEmoteController := crudcontrollers.NewNpcEmoteController(databaseResolver, logger)
	npcFactionController := crudcontrollers.NewNpcFactionController(databaseResolver, logger)
	npcFactionEntryController := crudcontrollers.NewNpcFactionEntryController(databaseResolver, logger)
	npcScaleGlobalBaseController := crudcontrollers.NewNpcScaleGlobalBaseController(databaseResolver, logger)
	npcSpellController := crudcontrollers.NewNpcSpellController(databaseResolver, logger)
	npcSpellsEffectController := crudcontrollers.NewNpcSpellsEffectController(databaseResolver, logger)
	npcSpellsEffectsEntryController := crudcontrollers.NewNpcSpellsEffectsEntryController(databaseResolver, logger)
	npcSpellsEntryController := crudcontrollers.NewNpcSpellsEntryController(databaseResolver, logger)
	npcTypeController := crudcontrollers.NewNpcTypeController(databaseResolver, logger)
	npcTypesTintController := crudcontrollers.NewNpcTypesTintController(databaseResolver, logger)
	objectContentController := crudcontrollers.NewObjectContentController(databaseResolver, logger)
	objectController := crudcontrollers.NewObjectController(databaseResolver, logger)
	perlEventExportSettingController := crudcontrollers.NewPerlEventExportSettingController(databaseResolver, logger)
	petController := crudcontrollers.NewPetController(databaseResolver, logger)
	petitionController := crudcontrollers.NewPetitionController(databaseResolver, logger)
	petsBeastlordDatumController := crudcontrollers.NewPetsBeastlordDatumController(databaseResolver, logger)
	petsEquipmentsetController := crudcontrollers.NewPetsEquipmentsetController(databaseResolver, logger)
	petsEquipmentsetEntryController := crudcontrollers.NewPetsEquipmentsetEntryController(databaseResolver, logger)
	playerTitlesetController := crudcontrollers.NewPlayerTitlesetController(databaseResolver, logger)
	questGlobalController := crudcontrollers.NewQuestGlobalController(databaseResolver, logger)
	raidDetailController := crudcontrollers.NewRaidDetailController(databaseResolver, logger)
	raidMemberController := crudcontrollers.NewRaidMemberController(databaseResolver, logger)
	reportController := crudcontrollers.NewReportController(databaseResolver, logger)
	respawnTimeController := crudcontrollers.NewRespawnTimeController(databaseResolver, logger)
	ruleSetController := crudcontrollers.NewRuleSetController(databaseResolver, logger)
	ruleValueController := crudcontrollers.NewRuleValueController(databaseResolver, logger)
	saylinkController := crudcontrollers.NewSaylinkController(databaseResolver, logger)
	serverScheduledEventController := crudcontrollers.NewServerScheduledEventController(databaseResolver, logger)
	sharedTaskActivityStateController := crudcontrollers.NewSharedTaskActivityStateController(databaseResolver, logger)
	sharedTaskController := crudcontrollers.NewSharedTaskController(databaseResolver, logger)
	sharedTaskDynamicZoneController := crudcontrollers.NewSharedTaskDynamicZoneController(databaseResolver, logger)
	sharedTaskMemberController := crudcontrollers.NewSharedTaskMemberController(databaseResolver, logger)
	skillCapController := crudcontrollers.NewSkillCapController(databaseResolver, logger)
	spawn2Controller := crudcontrollers.NewSpawn2Controller(databaseResolver, logger)
	spawnConditionController := crudcontrollers.NewSpawnConditionController(databaseResolver, logger)
	spawnConditionValueController := crudcontrollers.NewSpawnConditionValueController(databaseResolver, logger)
	spawnEventController := crudcontrollers.NewSpawnEventController(databaseResolver, logger)
	spawnentryController := crudcontrollers.NewSpawnentryController(databaseResolver, logger)
	spawngroupController := crudcontrollers.NewSpawngroupController(databaseResolver, logger)
	spellBucketController := crudcontrollers.NewSpellBucketController(databaseResolver, logger)
	spellGlobalController := crudcontrollers.NewSpellGlobalController(databaseResolver, logger)
	spellsNewController := crudcontrollers.NewSpellsNewController(databaseResolver, logger)
	startZoneController := crudcontrollers.NewStartZoneController(databaseResolver, logger)
	startingItemController := crudcontrollers.NewStartingItemController(databaseResolver, logger)
	taskActivityController := crudcontrollers.NewTaskActivityController(databaseResolver, logger)
	taskController := crudcontrollers.NewTaskController(databaseResolver, logger)
	tasksetController := crudcontrollers.NewTasksetController(databaseResolver, logger)
	timerController := crudcontrollers.NewTimerController(databaseResolver, logger)
	titleController := crudcontrollers.NewTitleController(databaseResolver, logger)
	traderController := crudcontrollers.NewTraderController(databaseResolver, logger)
	tradeskillRecipeController := crudcontrollers.NewTradeskillRecipeController(databaseResolver, logger)
	tradeskillRecipeEntryController := crudcontrollers.NewTradeskillRecipeEntryController(databaseResolver, logger)
	trapController := crudcontrollers.NewTrapController(databaseResolver, logger)
	tributeController := crudcontrollers.NewTributeController(databaseResolver, logger)
	tributeLevelController := crudcontrollers.NewTributeLevelController(databaseResolver, logger)
	veteranRewardTemplateController := crudcontrollers.NewVeteranRewardTemplateController(databaseResolver, logger)
	zoneController := crudcontrollers.NewZoneController(databaseResolver, logger)
	zoneFlagController := crudcontrollers.NewZoneFlagController(databaseResolver, logger)
	zonePointController := crudcontrollers.NewZonePointController(databaseResolver, logger)
	bootCrudControllers := provideCrudControllers(aaAbilityController, aaRankController, aaRankEffectController, aaRankPrereqController, accountController, accountFlagController, accountIpController, accountRewardController, adventureDetailController, adventureMemberController, adventureStatController, adventureTemplateController, adventureTemplateEntryController, adventureTemplateEntryFlavorController, alternateCurrencyController, auraController, baseDatumController, blockedSpellController, bookController, botDatumController, bugController, bugReportController, buyerController, charCreateCombinationController, charCreatePointAllocationController, charRecipeListController, characterActivityController, characterAltCurrencyController, characterAlternateAbilityController, characterAuraController, characterBandolierController, characterBindController, characterBuffController, characterCorpseController, characterCorpseItemController, characterCurrencyController, characterDatumController, characterDisciplineController, characterEnabledtaskController, characterExpModifierController, characterExpeditionLockoutController, characterInspectMessageController, characterInstanceSafereturnController, characterItemRecastController, characterLanguageController, characterLeadershipAbilityController, characterMaterialController, characterMemmedSpellController, characterPeqzoneFlagController, characterPetBuffController, characterPetInfoController, characterPetInventoryController, characterPotionbeltController, characterSkillController, characterSpellController, characterTaskController, characterTaskTimerController, completedSharedTaskActivityStateController, completedSharedTaskController, completedSharedTaskMemberController, completedTaskController, contentFlagController, damageshieldtypeController, dataBucketController, dbStrController, discordWebhookController, discoveredItemController, doorController, dynamicZoneController, dynamicZoneMemberController, dynamicZoneTemplateController, eventlogController, expeditionController, expeditionLockoutController, expeditionMemberController, factionAssociationController, factionBaseDatumController, factionListController, factionListModController, factionValueController, fishingController, forageController, friendController, globalLootController, gmIpController, graveyardController, gridController, gridEntryController, groundSpawnController, groupIdController, guildController, guildMemberController, guildRankController, guildRelationController, hackerController, horseController, instanceListController, instanceListPlayerController, inventoryController, inventorySnapshotController, ipExemptionController, itemController, itemTickController, ldonTrapEntryController, ldonTrapTemplateController, levelExpModController, lfguildController, loginAccountController, loginApiTokenController, loginServerAdminController, loginServerListTypeController, loginWorldServerController, logsysCategoryController, lootdropController, lootdropEntryController, loottableController, loottableEntryController, mailController, merchantlistController, merchantlistTempController, nameFilterController, npcEmoteController, npcFactionController, npcFactionEntryController, npcScaleGlobalBaseController, npcSpellController, npcSpellsEffectController, npcSpellsEffectsEntryController, npcSpellsEntryController, npcTypeController, npcTypesTintController, objectContentController, objectController, perlEventExportSettingController, petController, petitionController, petsBeastlordDatumController, petsEquipmentsetController, petsEquipmentsetEntryController, playerTitlesetController, questGlobalController, raidDetailController, raidMemberController, reportController, respawnTimeController, ruleSetController, ruleValueController, saylinkController, serverScheduledEventController, sharedTaskActivityStateController, sharedTaskController, sharedTaskDynamicZoneController, sharedTaskMemberController, skillCapController, spawn2Controller, spawnConditionController, spawnConditionValueController, spawnEventController, spawnentryController, spawngroupController, spellBucketController, spellGlobalController, spellsNewController, startZoneController, startingItemController, taskActivityController, taskController, tasksetController, timerController, titleController, traderController, tradeskillRecipeController, tradeskillRecipeEntryController, trapController, tributeController, tributeLevelController, veteranRewardTemplateController, zoneController, zoneFlagController, zonePointController)
	userContextMiddleware := middleware.NewUserContextMiddleware(databaseResolver, cache, logger)
	readOnlyMiddleware := middleware.NewReadOnlyMiddleware(databaseResolver, logger)
	requestLogMiddleware := middleware.NewRequestLogMiddleware(client)
	router := NewRouter(bootAppControllerGroups, bootCrudControllers, userContextMiddleware, readOnlyMiddleware, requestLogMiddleware, cache)
	httpServeCommand := cmd.NewHttpServeCommand(logger, router)
	routesListCommand := cmd.NewRoutesListCommand(router, logger)
	generateConfigurationCommand := cmd.NewGenerateConfigurationCommand(databaseResolver, logger)
	spireMigrateCommand := cmd.NewSpireMigrateCommand(connections, logger)
	questApiParseCommand := cmd.NewQuestApiParseCommand(logger, parseService)
	questExampleTestCommand := cmd.NewQuestExampleTestCommand(logger, questExamplesGithubSourcer)
	generateRaceModelMapsCommand := cmd.NewGenerateRaceModelMapsCommand(logger)
	v := ProvideCommands(helloWorldCommand, generateModelsCommand, generateControllersCommand, httpServeCommand, routesListCommand, generateConfigurationCommand, spireMigrateCommand, questApiParseCommand, questExampleTestCommand, generateRaceModelMapsCommand)
	webBoot := desktop.NewWebBoot(logger, router)
	app := NewApplication(db, logger, cache, v, databaseResolver, connections, router, webBoot)
	return app, nil
}
