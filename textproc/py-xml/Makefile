# $NetBSD: Makefile,v 1.32 2019/04/26 13:14:18 maya Exp $

DISTNAME=	PyXML-0.8.4
PKGNAME=	${PYPKGPREFIX}-xml-0.8.4
PKGREVISION=	5
CATEGORIES=	textproc python
MASTER_SITES=	${MASTER_SITE_SOURCEFORGE:=pyxml/}

MAINTAINER=	pkgsrc-users@NetBSD.org
#HOMEPAGE=	http://pyxml.sourceforge.net/
COMMENT=	Collection of libraries to process XML with Python

CONFLICTS+=	PyXML-[0-9]*

BUILDLINK_API_DEPENDS.expat+=	expat>=1.95.6
PYSETUPBUILDARGS+=	--with-libexpat=${BUILDLINK_DIR}
PYSETUPBUILDARGS+=	--ldflags="${EXPAT_LDFLAGS}"
EXPAT_LDFLAGS=		${COMPILER_RPATH_FLAG}${BUILDLINK_PREFIX.expat}/lib
FIX_RPATH+=		EXPAT_LDFLAGS

PYTHON_VERSIONS_INCOMPATIBLE=	 36 37 # dead upstream, no plans of porting

PYXML_SCRIPTS=	scripts/xmlproc_parse \
		scripts/xmlproc_val

REPLACE_PYTHON+=	xml/dom/ext/c14n.py
REPLACE_PYTHON+=	xml/dom/html/GenerateHtml.py

post-extract:
.for s in ${PYXML_SCRIPTS}
	${MV} ${WRKSRC}/${s} ${WRKSRC}/${s}${PYVERSSUFFIX}
.endfor

REPLACE_PYTHON=	${PYXML_SCRIPTS:S/$/${PYVERSSUFFIX}/g}

.include "../../textproc/expat/buildlink3.mk"
.include "../../lang/python/distutils.mk"
.include "../../lang/python/application.mk"
.include "../../mk/bsd.pkg.mk"
