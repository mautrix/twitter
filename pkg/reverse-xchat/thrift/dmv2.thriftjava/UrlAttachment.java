package com.x.dmv2.thriftjava;

import android.gov.nist.core.C0003b;
import android.gov.nist.core.Separators;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import java.io.IOException;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000:\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\u000e\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0006\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\f\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0007\b\u0086\b\u0018\u0000 $2\u00020\u0001:\u0002%$B9\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0005\u001a\u0004\u0018\u00010\u0004\u0012\b\u0010\u0006\u001a\u0004\u0018\u00010\u0004\u0012\b\u0010\u0007\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\b\u001a\u0004\u0018\u00010\u0002¢\u0006\u0004\b\t\u0010\nJ\u0017\u0010\u000e\u001a\u00020\r2\u0006\u0010\f\u001a\u00020\u000bH\u0016¢\u0006\u0004\b\u000e\u0010\u000fJ\u0012\u0010\u0010\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0010\u0010\u0011J\u0012\u0010\u0012\u001a\u0004\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u0012\u0010\u0013J\u0012\u0010\u0014\u001a\u0004\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u0014\u0010\u0013J\u0012\u0010\u0015\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0015\u0010\u0011J\u0012\u0010\u0016\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0016\u0010\u0011JL\u0010\u0017\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0005\u001a\u0004\u0018\u00010\u00042\n\b\u0002\u0010\u0006\u001a\u0004\u0018\u00010\u00042\n\b\u0002\u0010\u0007\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\b\u001a\u0004\u0018\u00010\u0002HÆ\u0001¢\u0006\u0004\b\u0017\u0010\u0018J\u0010\u0010\u0019\u001a\u00020\u0002HÖ\u0001¢\u0006\u0004\b\u0019\u0010\u0011J\u0010\u0010\u001b\u001a\u00020\u001aHÖ\u0001¢\u0006\u0004\b\u001b\u0010\u001cJ\u001a\u0010 \u001a\u00020\u001f2\b\u0010\u001e\u001a\u0004\u0018\u00010\u001dHÖ\u0003¢\u0006\u0004\b \u0010!R\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010\"R\u0016\u0010\u0005\u001a\u0004\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0005\u0010#R\u0016\u0010\u0006\u001a\u0004\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0006\u0010#R\u0016\u0010\u0007\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0007\u0010\"R\u0016\u0010\b\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\b\u0010\"¨\u0006&"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/UrlAttachment;", "Lcom/bendb/thrifty/a;", "", "url", "Lcom/x/dmv2/thriftjava/UrlAttachmentImage;", "banner_image_media_hash_key", "favicon_image_media_hash_key", "display_title", "attachment_id", "<init>", "(Ljava/lang/String;Lcom/x/dmv2/thriftjava/UrlAttachmentImage;Lcom/x/dmv2/thriftjava/UrlAttachmentImage;Ljava/lang/String;Ljava/lang/String;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/String;", "component2", "()Lcom/x/dmv2/thriftjava/UrlAttachmentImage;", "component3", "component4", "component5", "copy", "(Ljava/lang/String;Lcom/x/dmv2/thriftjava/UrlAttachmentImage;Lcom/x/dmv2/thriftjava/UrlAttachmentImage;Ljava/lang/String;Ljava/lang/String;)Lcom/x/dmv2/thriftjava/UrlAttachment;", "toString", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/String;", "Lcom/x/dmv2/thriftjava/UrlAttachmentImage;", "Companion", "UrlAttachmentAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class UrlAttachment implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final String attachment_id;

    @JvmField
    @InterfaceC88465b
    public final UrlAttachmentImage banner_image_media_hash_key;

    @JvmField
    @InterfaceC88465b
    public final String display_title;

    @JvmField
    @InterfaceC88465b
    public final UrlAttachmentImage favicon_image_media_hash_key;

    @JvmField
    @InterfaceC88465b
    public final String url;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new UrlAttachmentAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/UrlAttachment$UrlAttachmentAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/UrlAttachment;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/UrlAttachment;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/UrlAttachment;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class UrlAttachmentAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public UrlAttachment m85956read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            String string = null;
            UrlAttachmentImage urlAttachmentImage = null;
            UrlAttachmentImage urlAttachmentImage2 = null;
            String string2 = null;
            String string3 = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new UrlAttachment(string, urlAttachmentImage, urlAttachmentImage2, string2, string3);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        if (s != 3) {
                            if (s != 4) {
                                if (s != 5) {
                                    C11272a.m14141a(protocol, b);
                                } else if (b == 11) {
                                    string3 = protocol.readString();
                                } else {
                                    C11272a.m14141a(protocol, b);
                                }
                            } else if (b == 11) {
                                string2 = protocol.readString();
                            } else {
                                C11272a.m14141a(protocol, b);
                            }
                        } else if (b == 12) {
                            urlAttachmentImage2 = (UrlAttachmentImage) UrlAttachmentImage.ADAPTER.read(protocol);
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    } else if (b == 12) {
                        urlAttachmentImage = (UrlAttachmentImage) UrlAttachmentImage.ADAPTER.read(protocol);
                    } else {
                        C11272a.m14141a(protocol, b);
                    }
                } else if (b == 11) {
                    string = protocol.readString();
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a UrlAttachment struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("UrlAttachment");
            if (struct.url != null) {
                protocol.mo14136v3("url", 1, (byte) 11);
                protocol.mo14137w0(struct.url);
            }
            if (struct.banner_image_media_hash_key != null) {
                protocol.mo14136v3("banner_image_media_hash_key", 2, (byte) 12);
                UrlAttachmentImage.ADAPTER.write(protocol, struct.banner_image_media_hash_key);
            }
            if (struct.favicon_image_media_hash_key != null) {
                protocol.mo14136v3("favicon_image_media_hash_key", 3, (byte) 12);
                UrlAttachmentImage.ADAPTER.write(protocol, struct.favicon_image_media_hash_key);
            }
            if (struct.display_title != null) {
                protocol.mo14136v3("display_title", 4, (byte) 11);
                protocol.mo14137w0(struct.display_title);
            }
            if (struct.attachment_id != null) {
                protocol.mo14136v3("attachment_id", 5, (byte) 11);
                protocol.mo14137w0(struct.attachment_id);
            }
            protocol.mo14134i0();
        }
    }

    public UrlAttachment(@InterfaceC88465b String str, @InterfaceC88465b UrlAttachmentImage urlAttachmentImage, @InterfaceC88465b UrlAttachmentImage urlAttachmentImage2, @InterfaceC88465b String str2, @InterfaceC88465b String str3) {
        this.url = str;
        this.banner_image_media_hash_key = urlAttachmentImage;
        this.favicon_image_media_hash_key = urlAttachmentImage2;
        this.display_title = str2;
        this.attachment_id = str3;
    }

    public static /* synthetic */ UrlAttachment copy$default(UrlAttachment urlAttachment, String str, UrlAttachmentImage urlAttachmentImage, UrlAttachmentImage urlAttachmentImage2, String str2, String str3, int i, Object obj) {
        if ((i & 1) != 0) {
            str = urlAttachment.url;
        }
        if ((i & 2) != 0) {
            urlAttachmentImage = urlAttachment.banner_image_media_hash_key;
        }
        UrlAttachmentImage urlAttachmentImage3 = urlAttachmentImage;
        if ((i & 4) != 0) {
            urlAttachmentImage2 = urlAttachment.favicon_image_media_hash_key;
        }
        UrlAttachmentImage urlAttachmentImage4 = urlAttachmentImage2;
        if ((i & 8) != 0) {
            str2 = urlAttachment.display_title;
        }
        String str4 = str2;
        if ((i & 16) != 0) {
            str3 = urlAttachment.attachment_id;
        }
        return urlAttachment.copy(str, urlAttachmentImage3, urlAttachmentImage4, str4, str3);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final String getUrl() {
        return this.url;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final UrlAttachmentImage getBanner_image_media_hash_key() {
        return this.banner_image_media_hash_key;
    }

    @InterfaceC88465b
    /* renamed from: component3, reason: from getter */
    public final UrlAttachmentImage getFavicon_image_media_hash_key() {
        return this.favicon_image_media_hash_key;
    }

    @InterfaceC88465b
    /* renamed from: component4, reason: from getter */
    public final String getDisplay_title() {
        return this.display_title;
    }

    @InterfaceC88465b
    /* renamed from: component5, reason: from getter */
    public final String getAttachment_id() {
        return this.attachment_id;
    }

    @InterfaceC88464a
    public final UrlAttachment copy(@InterfaceC88465b String url, @InterfaceC88465b UrlAttachmentImage banner_image_media_hash_key, @InterfaceC88465b UrlAttachmentImage favicon_image_media_hash_key, @InterfaceC88465b String display_title, @InterfaceC88465b String attachment_id) {
        return new UrlAttachment(url, banner_image_media_hash_key, favicon_image_media_hash_key, display_title, attachment_id);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof UrlAttachment)) {
            return false;
        }
        UrlAttachment urlAttachment = (UrlAttachment) other;
        return Intrinsics.m65267c(this.url, urlAttachment.url) && Intrinsics.m65267c(this.banner_image_media_hash_key, urlAttachment.banner_image_media_hash_key) && Intrinsics.m65267c(this.favicon_image_media_hash_key, urlAttachment.favicon_image_media_hash_key) && Intrinsics.m65267c(this.display_title, urlAttachment.display_title) && Intrinsics.m65267c(this.attachment_id, urlAttachment.attachment_id);
    }

    public int hashCode() {
        String str = this.url;
        int iHashCode = (str == null ? 0 : str.hashCode()) * 31;
        UrlAttachmentImage urlAttachmentImage = this.banner_image_media_hash_key;
        int iHashCode2 = (iHashCode + (urlAttachmentImage == null ? 0 : urlAttachmentImage.hashCode())) * 31;
        UrlAttachmentImage urlAttachmentImage2 = this.favicon_image_media_hash_key;
        int iHashCode3 = (iHashCode2 + (urlAttachmentImage2 == null ? 0 : urlAttachmentImage2.hashCode())) * 31;
        String str2 = this.display_title;
        int iHashCode4 = (iHashCode3 + (str2 == null ? 0 : str2.hashCode())) * 31;
        String str3 = this.attachment_id;
        return iHashCode4 + (str3 != null ? str3.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        String str = this.url;
        UrlAttachmentImage urlAttachmentImage = this.banner_image_media_hash_key;
        UrlAttachmentImage urlAttachmentImage2 = this.favicon_image_media_hash_key;
        String str2 = this.display_title;
        String str3 = this.attachment_id;
        StringBuilder sb = new StringBuilder("UrlAttachment(url=");
        sb.append(str);
        sb.append(", banner_image_media_hash_key=");
        sb.append(urlAttachmentImage);
        sb.append(", favicon_image_media_hash_key=");
        sb.append(urlAttachmentImage2);
        sb.append(", display_title=");
        sb.append(str2);
        sb.append(", attachment_id=");
        return C0003b.m4b(sb, str3, Separators.RPAREN);
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}