# $NetBSD: Makefile,v 1.32 2018/12/09 21:05:36 adam Exp $

DISTNAME=		Jalbum${VERSION}
PKGNAME=		jalbum-${VERSION}
CATEGORIES=		www
MASTER_SITES=		http://jalbum.net/download/
DIST_SUBDIR=		${PKGNAME_NOREV}
EXTRACT_SUFX=		.zip

MAINTAINER=		pkgsrc-users@NetBSD.org
HOMEPAGE=		http://jalbum.net/
COMMENT=		Web photo album generator
LICENSE=		gnu-lgpl-v2

VERSION=		8.5.3

USE_JAVA2=		1.5
PKG_JVMS_ACCEPTED=	jdk15 sun-jdk6

WRKSRC=			${WRKDIR}/Jalbum

INSTALLATION_DIRS+=	bin ${PREFIX}/share/jalbum

USE_TOOLS+=		pax

CHECK_SHLIBS_SKIP+=	share/jalbum/lib/*
CHECK_SSP_SKIP+=	${CHECK_SHLIBS_SKIP}

post-extract:
	${ECHO} "===> Removing invalid skins"
	(cd ${WRKSRC} && rm -rf skins/Chameleon/styles/Black*)

do-build:
	sed	-e s#@PKG_JAVA_HOME@#${PKG_JAVA_HOME}#g		\
		-e s#@PREFIX@#${PREFIX}#g			\
		< ${FILESDIR}/jalbum > ${WRKDIR}/jalbum

do-install:
	${INSTALL_SCRIPT} ${WRKDIR}/jalbum ${DESTDIR}${PREFIX}/bin
	cd ${WRKDIR}/Jalbum && pax -rw . ${DESTDIR}${PREFIX}/share/jalbum

.include "../../mk/java-vm.mk"
.include "../../mk/bsd.pkg.mk"
